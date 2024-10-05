package main

import (
	"github.com/minhnghia2k3/GOssage/internal/auth"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"github.com/minhnghia2k3/GOssage/internal/store/cache"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockUserStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		config:        cfg,
		storage:       mockStore,
		cacheStorage:  &mockCacheStore,
		logger:        logger,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
