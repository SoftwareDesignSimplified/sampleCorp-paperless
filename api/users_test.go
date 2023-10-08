package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/qwetu_petro/backend/db/mock"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/utils"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       utils.RandomString(10),
		HashedPassword: hashedPassword,
		Email:          utils.RandomEmail(),
		FullName:       utils.RandomString(10),
		Department:     utils.RandomDepartment(),
	}
	return
}

func requireBodyMatcher(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Empty(t, gotUser.HashedPassword)

}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := utils.CheckPasswordHash(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
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
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"full_name":  user.FullName,
				"department": user.Department,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
					FullName:       user.FullName,
					Department:     user.Department,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)

			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
				requireBodyMatcher(t, recorder.Body, user)
			},
		}, {
			name: "InternalError",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"full_name":  user.FullName,
				"department": user.Department,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 500, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"full_name":  user.FullName,
				"department": user.Department,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})

			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":   "invalid-user#name",
				"password":   password,
				"email":      user.Email,
				"full_name":  user.FullName,
				"department": user.Department,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":   user.Username,
				"password":   "123",
				"email":      user.Email,
				"full_name":  user.FullName,
				"department": user.Department,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// {
		// 	name: "wrongPassword",
		// 	body: gin.H{
		// 		"username":  user.Username,
		// 		"password":  "123",
		// 		"email":     user.Email,
		// 		"full_name": user.FullName,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)

		// 	},
		// 	checkRes: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, 400, recorder.Code)
		// 	},
		// }, {
		// 	name: "StatusBadRequest",
		// 	body: gin.H{
		// 		"username":  user.Username,
		// 		"password":  "123",
		// 		"email":     user.Email,
		// 		"full_name": user.FullName,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkRes: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/auth/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(recorder)
		})
	}

}
