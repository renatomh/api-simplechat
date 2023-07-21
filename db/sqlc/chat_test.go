package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateChat(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 3; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contacts between users
	contact1, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	contact2, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[2].ID,
	})
	testQueries.AcceptContact(context.Background(), contact1.ID)
	testQueries.AcceptContact(context.Background(), contact2.ID)

	// Creating chats for the contacts
	chat1, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact1.FromUserID,
		ToUserID:   contact1.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat1)
	chat2, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact2.FromUserID,
		ToUserID:   contact2.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat2)
}

func TestGetChat(t *testing.T) {
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
	// Creating chat for the contacts
	chat, _ := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})

	// Querying the newly created chat by its ID
	queriedChat, err := testQueries.GetChat(context.Background(), chat.ID)

	require.NoError(t, err, "unexpected error querying the chat: %v", err)
	require.NotEmpty(t, queriedChat)

	require.Equal(t, chat.ID, queriedChat.ID)
	require.Equal(t, chat.FromUserID, queriedChat.FromUserID)
	require.Equal(t, chat.ToUserID, queriedChat.ToUserID)
	require.Empty(t, chat.LastMessageReceivedAt)
}

func TestListChats(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 3; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contacts between users
	contact1, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	contact2, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[2].ID,
	})
	testQueries.AcceptContact(context.Background(), contact1.ID)
	testQueries.AcceptContact(context.Background(), contact2.ID)

	// Creating chats for the contacts
	testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact1.FromUserID,
		ToUserID:   contact1.ToUserID,
	})
	testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact2.FromUserID,
		ToUserID:   contact2.ToUserID,
	})

	// Listing the chats
	chats, err := testQueries.ListChats(context.Background(),
		ListChatsParams{
			FromUserID: users[0].ID,
			Limit:      10,
			Offset:     0,
		})
	require.NoError(t, err)
	require.Len(t, chats, 2)
	for _, chat := range chats {
		require.NotEmpty(t, chat)
	}
}

func TestDeleteChat(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 3; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	testQueries.AcceptContact(context.Background(), contact.ID)

	// Creating chat for the contact
	chat, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})
	require.NoError(t, err)

	// Deleting the chat
	err = testQueries.DeleteChat(context.Background(), chat.ID)
	require.NoError(t, err)

	// Checking if chat was deleted
	deletedChat, err := testQueries.GetChat(context.Background(), chat.ID)
	require.Error(t, err)
	require.Empty(t, deletedChat)
}
