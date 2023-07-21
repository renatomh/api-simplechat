package db

import (
	"context"
	"testing"
	"time"

	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

func TestCreateContact(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 4; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating contacts between users
	contact1, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, contact1)
	contact2, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[2].ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, contact2)
	contact3, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[3].ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, contact3)
	contact4, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[1].ID,
		ToUserID:   users[2].ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, contact4)
}

func TestGetContact(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating contact between users
	contact, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})

	// Querying the newly created contact by its ID
	queriedContact, err := testQueries.GetContact(context.Background(), contact.ID)

	require.NoError(t, err, "unexpected error querying the contact: %v", err)
	require.NotEmpty(t, queriedContact)

	require.Equal(t, contact.ID, queriedContact.ID)
	require.Equal(t, contact.FromUserID, queriedContact.FromUserID)
	require.Equal(t, contact.ToUserID, queriedContact.ToUserID)
	require.Equal(t, contact.Status, queriedContact.Status)
	require.WithinDuration(t, contact.RequestedAt, queriedContact.RequestedAt, time.Second)
}

func TestAcceptContact(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating contact between users
	contact, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})

	// Accepting the requested contact
	acceptedContact, err := testQueries.AcceptContact(context.Background(), contact.ID)

	require.NoError(t, err, "unexpected error accepting the contact: %v", err)
	require.NotEmpty(t, acceptedContact)

	require.Equal(t, contact.ID, acceptedContact.ID)
	require.Equal(t, contact.FromUserID, acceptedContact.FromUserID)
	require.Equal(t, contact.ToUserID, acceptedContact.ToUserID)
	require.Equal(t, "Accepted", acceptedContact.Status)
	require.WithinDuration(t, contact.RequestedAt, acceptedContact.RequestedAt, time.Second)
	require.WithinDuration(t, time.Now(), acceptedContact.AcceptedAt.Time, time.Second)
}

func TestRejectContact(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating contact between users
	contact, err := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})

	// Rejecting the requested contact
	rejectedContact, err := testQueries.RejectContact(context.Background(), contact.ID)

	require.NoError(t, err, "unexpected error rejecting the contact: %v", err)
	require.NotEmpty(t, rejectedContact)

	require.Equal(t, contact.ID, rejectedContact.ID)
	require.Equal(t, contact.FromUserID, rejectedContact.FromUserID)
	require.Equal(t, contact.ToUserID, rejectedContact.ToUserID)
	require.Equal(t, "Rejected", rejectedContact.Status)
	require.WithinDuration(t, contact.RequestedAt, rejectedContact.RequestedAt, time.Second)
	require.Empty(t, rejectedContact.AcceptedAt)
}

func TestListContacts(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 5; i++ {
		user, err := createRandomUser(t)
		require.NoError(t, err)
		users = append(users, user)
	}

	// Initializing counters for pending, accepted and rejected contacts
	pendingCount := 0
	acceptedCount := 0
	rejectedCount := 0
	// Creating contacts between users
	contacts := []Contact{}
	for i := 1; i < 5; i++ {
		contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
			FromUserID: users[0].ID,
			ToUserID:   users[i].ID,
		})
		statusId := util.RandomInt(0, 2)
		if statusId == 0 {
			pendingCount++
		} else if statusId == 1 {
			acceptedCount++
			contact, _ = testQueries.AcceptContact(context.Background(), contact.ID)
		} else if statusId == 2 {
			rejectedCount++
			contact, _ = testQueries.RejectContact(context.Background(), contact.ID)
		}
		contacts = append(contacts, contact)
	}

	// Now, we'll list the contacts
	allContacts, err := testQueries.ListContacts(context.Background(), ListContactsParams{
		FromUserID: users[0].ID,
		Limit:      10,
		Offset:     0,
	})
	require.NoError(t, err)
	require.Len(t, allContacts, pendingCount+acceptedCount+rejectedCount)

	// Then, we'll list the pending contacts
	pendingContacts, err := testQueries.ListPendingContacts(context.Background(),
		ListPendingContactsParams{
			FromUserID: users[0].ID,
			Limit:      10,
			Offset:     0,
		})
	require.NoError(t, err)
	require.Len(t, pendingContacts, pendingCount)

	// Listing the accepted contacts
	acceptedContacts, err := testQueries.ListAcceptedContacts(context.Background(),
		ListAcceptedContactsParams{
			FromUserID: users[0].ID,
			Limit:      10,
			Offset:     0,
		})
	require.NoError(t, err)
	require.Len(t, acceptedContacts, acceptedCount)
	for _, contact := range acceptedContacts {
		require.NotEmpty(t, contact)
		require.Equal(t, "Accepted", contact.Status)
		require.NotEmpty(t, contact.AcceptedAt)
	}

	// Listing the rejected contacts
	rejectedContacts, err := testQueries.ListRejectedContacts(context.Background(),
		ListRejectedContactsParams{
			FromUserID: users[0].ID,
			Limit:      10,
			Offset:     0,
		})
	require.NoError(t, err)
	require.Len(t, rejectedContacts, rejectedCount)
	for _, contact := range rejectedContacts {
		require.NotEmpty(t, contact)
		require.Equal(t, "Rejected", contact.Status)
		require.Empty(t, contact.AcceptedAt)
	}
}

func TestDeleteContact(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})

	// Deleting the contact
	err := testQueries.DeleteContact(context.Background(), contact.ID)
	require.NoError(t, err)

	// Checking if contact was deleted
	deletedContact, err := testQueries.GetContact(context.Background(), contact.ID)
	require.Error(t, err)
	require.Empty(t, deletedContact)
}
