package twitter_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jake-hansen/followrs/repositories/apis/twitter"
	"github.com/stretchr/testify/assert"
)

func NewTestServer() (*http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return mux, server
}

type emptyBody struct{}

func StandardHandler(t *testing.T, mux *http.ServeMux, endpoint string, body interface{}) {
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-rate-limit-remaining", "100")
		w.Header().Set("x-rate-limit-reset", "100")
		bytes, _ := json.Marshal(body)
		w.Write(bytes)
	})
}

func TestPerformRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		StandardHandler(t, mux, "/test", nil)
		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))

		assert.NoError(t, err)
		assert.Equal(t, int64(100), endpoint.RemainingCalls)
		assert.Equal(t, time.Unix(100, 0), endpoint.RateLimitReset)
	})

	t.Run("no-rate-limit-header-present", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-rate-limit-reset", "100")
			bytes, _ := json.Marshal(new(emptyBody))
			w.Write(bytes)
		})

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))

		assert.Error(t, err)
		assert.Equal(t, "header x-rate-limit-remaining not found in response", err.Error())
	})

	t.Run("no-rate-limit-reset-header-present", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-rate-limit-remaining", "100")
			bytes, _ := json.Marshal(new(emptyBody))
			w.Write(bytes)
		})

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))

		assert.Error(t, err)
		assert.Equal(t, "header x-rate-limit-reset not found in response", err.Error())
	})

	t.Run("error-parsing-rate-limit-header", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-rate-limit-remaining", "not a number!")
			w.Header().Set("x-rate-limit-reset", "100")
			bytes, _ := json.Marshal(new(emptyBody))
			w.Write(bytes)
		})

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse rate limit header x-rate-limit-remaining")
	})

	t.Run("error-parsing-rate-limit-reset-header", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-rate-limit-remaining", "100")
			w.Header().Set("x-rate-limit-reset", "not a number!")
			bytes, _ := json.Marshal(new(emptyBody))
			w.Write(bytes)
		})

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse rate limit header x-rate-limit-reset")
	})

	t.Run("rate-limit-reached", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		endpoint := &twitter.Endpoint{
			URL:            "/test",
			RemainingCalls: 0,
			RateLimitReset: time.Time{},
		}

		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			durationAddition, _ := time.ParseDuration("5s")
			resetTime := time.Now().Add(durationAddition).Unix()

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-rate-limit-remaining", "0")
			w.Header().Set("x-rate-limit-reset", strconv.Itoa(int(resetTime)))
			bytes, _ := json.Marshal(new(emptyBody))
			w.Write(bytes)
		})

		api, _ := twitter.NewTwitterAPI(server.URL, "", "", "")
		req, _ := retryablehttp.NewRequest("GET", "/test", nil)
		err := endpoint.PerformRequest(req, api.Client, new(emptyBody))
		assert.NoError(t, err)
		err = endpoint.PerformRequest(req, api.Client, new(emptyBody))
		assert.Error(t, err)
	})

}
