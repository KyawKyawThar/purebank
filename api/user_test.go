package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "purebank/db/mock"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	//fmt.Println("e.arg is ", e.arg)
	//fmt.Println("e.password is: ", e.password)

	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	//fmt.Println("arg is: ", arg)

	//plain pass and hash password
	err := util.CompareHashPassword(e.password, e.arg.Password)

	if err != nil {
		fmt.Println("CODE SHOULDN'T .......")
		return false
	}

	arg.Password = e.arg.Password

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {

	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {

	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {

	user, password := randomUser(t)

	testCases := []struct {
		status        string
		body          gin.H
		buildStubs    func(s *mockdb.MockStore)
		checkResponse func(r *httptest.ResponseRecorder)
	}{
		{
			status: "OK",
			body: gin.H{
				"username":   user.Username,
				"email":      user.Email,
				"password":   user.Password,
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:  user.Username,
					//Password:  user.Password,
					Email:     user.Email,
					FirstName: user.FirstName,
				}
				s.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {

				fmt.Println("Check response: ", r.Body)
				fmt.Println("Check User: ", user)

				requireBodyMatchUser(t, r.Body, user)
				require.Equal(t, http.StatusOK, r.Code)

			},
		},
		{
			status: "InternalError",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Users{}, sql.ErrConnDone)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
		{
			status: "Invalid Hash Password",
			body: gin.H{
				"username":   user.Username,
				"password":   util.RandomString(73),
				"email":      user.Email,
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
		{
			status: "DuplicatedUserName",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Users{}, &pgconn.PgError{Routine: "_bt_check_unique"})
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				fmt.Println("checkResponse:", r.Code)
				require.Equal(t, http.StatusForbidden, r.Code)
			},
		},

		{
			status: "TooShortPassword",
			body: gin.H{
				"username":   user.Username,
				"password":   "123",
				"email":      user.Email,
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			status: "InvalidEmail",
			body: gin.H{
				"username":   user.Username,
				"password":   "123",
				"email":      "invalid-email",
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			status: "InvalidUsername",
			body: gin.H{
				"username":   "invalid-user#1",
				"password":   password,
				"email":      "invalid-email",
				"first_name": user.FirstName,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.status, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			server := newTestServer(t, store)

			r := httptest.NewRecorder()

			fmt.Println("recorder in request: ", r)
			server.router.ServeHTTP(r, request)

			tc.checkResponse(r)
		})
	}
}

func TestLogInUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		status        string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(r *httptest.ResponseRecorder)
	}{
		{
			status: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)

			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
			},
		},

		{
			status: "InternalError",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"first_name": user.FirstName,
				"email":      user.Email,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).
					Return(db.Users{}, sql.ErrConnDone)
			},

			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},

		{
			status: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": "password",
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).
					Return(user, sql.ErrNoRows)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, r.Code)
			},
		},

		{
			status: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).
					Return(user, nil)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
		{
			status: "Invalid username",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for i := range testCases {
		fmt.Println("i is:", i)
		tc := testCases[i]

		t.Run(tc.status, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			r := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			url := "/user/login"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(r, req)

			tc.checkResponse(r)
		})
	}
}

// randomUser for create User API
func randomUser(t *testing.T) (user db.Users, password string) {

	password = util.RandomString(10)
	hashPassword, err := util.HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	user = db.Users{
		Username:  util.RandomOwner(6),
		Password:  hashPassword,
		Email:     util.RandomEmail(6),
		FirstName: util.RandomString(5),
	}
	return

}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.Users) {

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.Users

	err = json.Unmarshal(data, &gotUser)

	//fmt.Println("got User", gotUser.Password)
	//fmt.Println("User is:", user.Password)
	require.Equal(t, gotUser.Username, user.Username)
	require.Equal(t, user.FirstName, gotUser.FirstName)
	require.Equal(t, user.Email, gotUser.Email)

	require.Empty(t, gotUser.Password)
}
