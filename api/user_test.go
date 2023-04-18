package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockdb "purebank/db/mock"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"testing"
)

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
				fmt.Println("user is:", user)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				fmt.Println("r.body is", r.Result())
				require.Equal(t, http.StatusOK, r.Code)
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

			server, err := NewServer(store)

			require.NoError(t, err)

			r := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			fmt.Println("data is", tc.body)
			fmt.Println("err is ", err)
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



	//data, err := io.ReadAll(body)
	//require.NoError(t, err)
	//
	//var gotUser db.Users
	//
	//err = json.Unmarshal(data, &gotUser)
	//fmt.Println("got User", gotUser)
	//fmt.Println("User is:", user)
	//require.Equal(t, gotUser.Username, user.Username)
	//require.Equal(t, user.FirstName, gotUser.FirstName)
	//require.Equal(t, user.Email, gotUser.Email)
	//
	//require.Empty(t, user.Password, gotUser.Password)
}
