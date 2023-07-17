package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// Retrieving a random username
	username := util.RandomUsername()
	// Creating user with email
	arg := CreateUserParams{
		Name:     "John Doe",
		Username: username,
		Email:    sql.NullString{String: username + "@domain.com", Valid: true},
	}

	user1, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user1)

	require.Equal(t, arg.Name, user1.Name)
	require.Equal(t, arg.Username, user1.Username)
	require.Equal(t, arg.Email, user1.Email)

	require.NotZero(t, user1.ID)
	require.NotZero(t, user1.CreatedAt)

	// Trying to create another user user with the same username
	user2, err := testQueries.CreateUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, user2)

	// Creating user without email
	arg = CreateUserParams{
		Name:     "Jack Doe",
		Username: util.RandomUsername(),
		Email:    sql.NullString{},
	}

	user3, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user3)

	require.Equal(t, arg.Name, user3.Name)
	require.Equal(t, arg.Username, user3.Username)
	require.Equal(t, arg.Email, arg.Email)

	require.NotZero(t, user3.ID)
	require.NotZero(t, user3.CreatedAt)

	// At the end, we could delete the created users, to avoid filling the database
	testQueries.DeleteUser(context.Background(), user1.ID)
	// user2 is not created, so we don't need to remove it
	testQueries.DeleteUser(context.Background(), user3.ID)

}
