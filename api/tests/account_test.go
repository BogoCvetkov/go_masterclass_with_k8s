package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BogoCvetkov/go_mastercalss/auth"
	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	mockdb "github.com/BogoCvetkov/go_mastercalss/db/mock"
	testutil "github.com/BogoCvetkov/go_mastercalss/db/tests"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker interfaces.IAuth,
	authorizationType string,
	userId int64,
	duration time.Duration,
) {
	p1, _ := auth.NewTokenPayload(userId, time.Duration(time.Minute*15))
	token, err := tokenMaker.GenerateToken(p1)
	require.NoError(t, err)
	require.NotEmpty(t, p1)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("Authorization", authorizationHeader)
}

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.ID)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth)
		buildStubs    func(store *mockdb.MockIStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)
				store.EXPECT().
					GetAccountByOwner(gomock.Any(), gomock.Eq(db.GetAccountByOwnerParams{ID: account.ID, Owner: user.ID})).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", 928138123, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(int64(928138123))).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().
					GetAccountByOwner(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetAccountByOwner(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: int64(1241421),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},

			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)
				store.EXPECT().
					GetAccountByOwner(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					GetAccountByOwner(gomock.Any(), gomock.Eq(db.GetAccountByOwnerParams{ID: 0, Owner: user.ID})).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockIStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.GetAuth())
			server.GetRouter().ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth)
		buildStubs    func(store *mockdb.MockIStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
			},
			buildStubs: func(store *mockdb.MockIStore) {

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(0)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockIStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/account"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.GetAuth())
			server.GetRouter().ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.ID)
	}

	type Query struct {
		Page  int
		Limit int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth)
		buildStubs    func(store *mockdb.MockIStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Page:  1,
				Limit: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {
				arg := db.ListAccountsParams{
					Owner:  user.ID,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				Page:  1,
				Limit: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
			},
			buildStubs: func(store *mockdb.MockIStore) {
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(0)

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				Page:  -1,
				Limit: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				Page:  1,
				Limit: -1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker interfaces.IAuth) {
				addAuthorization(t, request, tokenMaker, "Bearer", user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockIStore) {

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(*user, nil)

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockIStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/account"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.Page))
			q.Add("limit", fmt.Sprintf("%d", tc.query.Limit))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.GetAuth())
			server.GetRouter().ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAccount(owner int64) db.Account {
	return db.Account{
		ID:       testutil.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  testutil.RandomMoney(),
		Currency: testutil.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
