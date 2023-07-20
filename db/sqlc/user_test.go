package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserViaAPI(t *testing.T) (User, error) {
	// Retrieving a random user from the third-party API
	randomUsers, err := util.QueryRandomUsersAPI(1)
	require.NoError(t, err, "unexpected error querying random users API: %v", err)
	arg := CreateUserParams{
		FullName: fmt.Sprintf("%s %s", randomUsers[0].Name.First, randomUsers[0].Name.Last),
		Username: randomUsers[0].Login.Username,
		Email:    sql.NullString{String: randomUsers[0].Email, Valid: true},
		HashPass: randomUsers[0].Login.Password,
	}

	// Creating user in the database
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashPass, user.HashPass)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	// Returning created user
	return user, nil
}

func createRandomUser(t *testing.T) (User, error) {
	// Retrieving a random user from the local functions
	username := util.RandomUsername()
	arg := CreateUserParams{
		FullName: fmt.Sprintf(
			"%s %s",
			strings.Title(strings.Split(username, ".")[0]),
			strings.Title(strings.Split(username, ".")[1]),
		),
		Username: username,
		Email: sql.NullString{
			String: username + "@" + util.RandomString(int(util.RandomInt(6, 9))) + ".com",
			Valid:  true,
		},
		HashPass: util.RandomString(int(util.RandomInt(18, 24))),
	}

	// Creating user in the database
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashPass, user.HashPass)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	// Returning created user
	return user, nil
}

func TestCreateUser(t *testing.T) {
	// Creating a random user
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// Creating a random user
	createdUser, err := createRandomUser(t)
	require.NoError(t, err, "unexpected error creating the user: %v", err)

	// Querying the newly created user by its ID
	queriedUser, err := testQueries.GetUser(context.Background(), createdUser.ID)

	require.NoError(t, err, "unexpected error querying the user: %v", err)
	require.NotEmpty(t, queriedUser)

	require.Equal(t, createdUser.ID, queriedUser.ID)
	require.Equal(t, createdUser.FullName, queriedUser.FullName)
	require.Equal(t, createdUser.Username, queriedUser.Username)
	require.Equal(t, createdUser.Email, queriedUser.Email)
	require.Equal(t, createdUser.HashPass, queriedUser.HashPass)
	require.WithinDuration(t, createdUser.CreatedAt, queriedUser.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	// Creating a random user
	createdUser, err := createRandomUser(t)
	require.NoError(t, err, "unexpected error creating the user: %v", err)

	args := UpdateUserParams{
		ID:        createdUser.ID,
		FullName:  "Mr." + createdUser.FullName,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		AvatarUrl: createdUser.AvatarUrl,
	}

	// Updating the newly created user
	updatedUser, err := testQueries.UpdateUser(context.Background(), args)

	require.NoError(t, err, "unexpected error updating the user: %v", err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, createdUser.ID, updatedUser.ID)
	require.Equal(t, args.FullName, updatedUser.FullName)
	require.Equal(t, createdUser.Username, updatedUser.Username)
	require.Equal(t, createdUser.Email, updatedUser.Email)
	require.Equal(t, createdUser.HashPass, updatedUser.HashPass)
	require.WithinDuration(t, createdUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	// Creating a random user
	createdUser, err := createRandomUser(t)
	require.NoError(t, err, "unexpected error creating the user: %v", err)

	// Deleting the user
	err = testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)

	// Checking if user was deleted
	queriedUser, err := testQueries.GetUser(context.Background(), createdUser.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, queriedUser)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	args := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestChangeUserPassword(t *testing.T) {
	// Creating a random user
	createdUser, err := createRandomUser(t)
	require.NoError(t, err, "unexpected error creating the user: %v", err)

	args := ChangeUserPasswordParams{
		ID:       createdUser.ID,
		HashPass: util.RandomString(int(util.RandomInt(18, 24))),
	}

	// Changing the newly created user's password
	updatedUser, err := testQueries.ChangeUserPassword(context.Background(), args)

	require.NoError(t, err, "unexpected error changing user's password: %v", err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, createdUser.ID, updatedUser.ID)
	require.Equal(t, args.HashPass, updatedUser.HashPass)
	require.WithinDuration(t, updatedUser.PasswordChangedAt, time.Now(), time.Second)
}
