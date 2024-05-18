package traefik_redirect_on_error

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectOnStatusPlugin_ServeHTTP(t *testing.T) {
	testCases := []struct {
		desc     string
		next     http.HandlerFunc
		ctx      context.Context
		config   *Config
		method   string
		validate func(t *testing.T, recorder *httptest.ResponseRecorder, config *Config)
	}{
		{
			desc: "Should redirect",
			next: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(502)
			}),
			ctx: context.Background(),
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 307,
				Status:       []string{"502"},
				Method:       nil,
			},
			method: http.MethodGet,
			validate: func(t *testing.T, recorder *httptest.ResponseRecorder, config *Config) {
				t.Helper()
				assertStatus(t, recorder, config.RedirectCode)
				assertHeader(t, recorder, "Location", config.RedirectUri)
			},
		},
		{
			desc: "Should not redirect",
			next: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte("Test"))
			}),
			ctx: context.Background(),
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 307,
				Status:       []string{"502"},
				Method:       nil,
			},
			method: http.MethodGet,
			validate: func(t *testing.T, recorder *httptest.ResponseRecorder, config *Config) {
				t.Helper()
				assertStatus(t, recorder, 200)
				if !bytes.Equal(recorder.Body.Bytes(), []byte("Test")) {
					t.Error("Body does not equal")
				}
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			redirectOnStatusHandler, err := New(test.ctx, test.next, test.config, "test")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(test.ctx, test.method, "http://localhost", nil)

			if err != nil {
				t.Fatal(err)
			}

			redirectOnStatusHandler.ServeHTTP(recorder, req)

			test.validate(t, recorder, test.config)
		})
	}
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if recorder.Code != expected {
		t.Errorf("Response Code should equal to redirect status code. Expected: %d, Actual: %d", expected, recorder.Code)
	}
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, expected string) {
	t.Helper()
	if recorder.Header().Get(key) != expected {
		t.Errorf("invalid header value for key %s: Expected: %s Actual: %s", key, expected, recorder.Header().Get(key))
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		desc      string
		config    *Config
		shouldErr bool
	}{
		{
			desc:      "Should not return error for redirect 307",
			shouldErr: false,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 307,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should not return error for redirect 302",
			shouldErr: false,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 302,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should not return error for redirect 303",
			shouldErr: false,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 303,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for redirect 304",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 304,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for redirect 200",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 200,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for redirect 400",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 400,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for empty redirect uri",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "",
				RedirectCode: 307,
				Status:       []string{"400-600"},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for empty status",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 307,
				Status:       []string{},
				Method:       nil,
			},
		},
		{
			desc:      "Should return error for invalid status",
			shouldErr: true,
			config: &Config{
				RedirectUri:  "/test",
				RedirectCode: 307,
				Status:       []string{"abc"},
				Method:       nil,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			_, err := New(ctx, handler, test.config, "test")
			if test.shouldErr && err == nil {
				t.Errorf("Did not return error for testcase: %s", test.desc)
			}

			if !test.shouldErr && err != nil {
				t.Errorf("Did return error for testcase: %s", test.desc)
			}
		})
	}
}
