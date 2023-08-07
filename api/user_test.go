package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/renatomh/api-simplechat/db/mock"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

// Struct to hold fields for custom matcher
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// Matches checks if provided args match, and returns the result
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// Checking if args are present in the interface
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// Checking if passwords match
	err := util.CheckPassword(e.password, arg.HashPass)
	if err != nil {
		return false
	}

	e.arg.HashPass = arg.HashPass
	return reflect.DeepEqual(e.arg, arg)
}

// String defines the matcher message
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// Matcher for the provided and hashed password
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Initializing the gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Initializing a new store with the gomock controller
			store := mockdb.NewMockStore(ctrl)

			// Building stubs for the test case
			tc.buildStubs(store)

			// Starting test server and sending request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// Defining request URL and making the request
			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// Here, we'll serve the requests and save it in the recorder
			server.router.ServeHTTP(recorder, request)

			// Checking the response
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	// Creating a random user
	user, _ := randomUser(t)

	// Defining tests cases
	testCases := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Building stubs to check the calling of GetUser method
				store.EXPECT().
					// Expects to have the same user ID
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					// Should be called only once
					Times(1).
					// And expects to return the user object and a nil error
					// The return params musth match the 'querier' function declaration
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Checking if response is correct
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Building stubs to check the calling of GetUser method
				store.EXPECT().
					// Expects to have the same user ID
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					// Should be called only once
					Times(1).
					// And expects to return an empty user object and an error
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Checking if response is correct
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Building stubs to check the calling of GetUser method
				store.EXPECT().
					// Expects to have the same user ID
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					// Should be called only once
					Times(1).
					// And expects to return an empty user object and an error
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Checking if response is correct
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// Building stubs to check the calling of GetUser method
				store.EXPECT().
					// Will use an invalid ID
					GetUser(gomock.Any(), gomock.Any()).
					// Should not be called
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Checking if response is correct
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// Testing the defined cases
	for i := range testCases {
		// Getting the test case to be checked
		tc := testCases[i]

		// Running the test case
		t.Run(tc.name, func(t *testing.T) {
			// Initializing the gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Initializing a new store with the gomock controller
			store := mockdb.NewMockStore(ctrl)

			// Building stubs for the test case
			tc.buildStubs(store)

			// Starting test server and sending request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Defining request URL and making the request
			url := fmt.Sprintf("/users/%d", tc.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Here, we'll serve the requests and save it in the recorder
			server.router.ServeHTTP(recorder, request)

			// Checking the response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	// Retrieving a random user from the local functions
	username := util.RandomUsername()
	password = util.RandomString(6)
	hashpass, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
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

	// Returning both user and password (as defined in the function signature)
	return
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
