package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/renatomh/api-simplechat/db/mock"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

func TestGetUserAPI(t *testing.T) {
	// Creating a random user
	user := randomUser()

	// Initializing the gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initializing a new store with the gomock controller
	store := mockdb.NewMockStore(ctrl)

	// Building stubs to check the calling of GetUser method
	store.EXPECT().
		// Expects to have the same user ID
		GetUser(gomock.Any(), gomock.Eq(user.ID)).
		// Should be called only once
		Times(1).
		// And expects to return the user object and a nil error
		// The return params musth match the 'querier' function declaration
		Return(user, nil)

	// Starting test server and sending request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	// Defining request URL and making the request
	url := fmt.Sprintf("/users/%d", user.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Here, we'll serve the requests and save it in the recorder
	server.router.ServeHTTP(recorder, request)

	// Checking if response is correct
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchUser(t, recorder.Body, user)
}

func randomUser() db.User {
	// Retrieving a random user from the local functions
	username := util.RandomUsername()
	hashpass, _ := util.HashPassword(util.RandomString(int(util.RandomInt(18, 24))))
	return db.User{
		ID: util.RandomInt(1, 1000),
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
		HashPass: hashpass,
	}
}

// requireBodyMatchUser checks if the response body for the request is correct
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	// Reading data from the response body
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// We'll unmarshal only the fields returned by the route
	var gotUser publicUserResponse

	// If we tried to unmarshal using the 'db.User' struct, it would lead to errors
	// Since we can't unmarshal 'string' to 'sql.NullString' and others
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email.String, gotUser.Email)
	require.Equal(t, user.AvatarUrl.String, gotUser.AvatarUrl)
	require.Equal(t, user.LastLoginAt.Time, gotUser.LastLoginAt)
}
