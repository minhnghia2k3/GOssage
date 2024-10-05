package main

import (
	"github.com/minhnghia2k3/GOssage/internal/store/cache"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	// > Declare with config with redisConfig
	withRedis := config{
		redisConfig: redisConfig{
			enabled: true,
		},
	}

	app := newTestApplication(t, withRedis)
	mux := app.mount()

	testToken, _ := app.authenticator.GenerateToken(nil)

	// > Test case #1: not allow unauthenticated request.
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	// > Test case #2: Should allow authenticated request.
	t.Run("should allow authenticated requests", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		mockCacheStore.On("Get", int64(1)).Return(nil, nil).Twice()

		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.Calls = nil
	})

	// > Test case #3: Should hit the cache first and if not exists
	// set the user on the cache
	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockCacheStore.On("Get", int64(42)).Return(nil, nil)
		mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.Calls = nil // Reset mock expectations
	})

	// > Test case #4: Should Not hit the cache if it is not enabled

	t.Run("should NOT hit the cache if it is not enabled", func(t *testing.T) {
		withRedis := config{
			redisConfig: redisConfig{
				enabled: false,
			},
		}
		app := newTestApplication(t, withRedis)
		mux := app.mount()
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNotCalled(t, "Get")
		mockCacheStore.Calls = nil // Reset mock expectations
	})

}
