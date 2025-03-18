package shared

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFormRequest(t *testing.T) {
	tests := []struct {
		name       string
		data       url.Values
		statusCode int
		response   string
		expectErr  bool
	}{
		{
			name: "successful post",
			data: url.Values{
				"key1": []string{"value1"},
				"key2": []string{"value2"},
			},
			statusCode: http.StatusOK,
			response:   "success response",
			expectErr:  false,
		},
		{
			name: "server error",
			data: url.Values{
				"key": []string{"value"},
			},
			statusCode: http.StatusInternalServerError,
			response:   "error response",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			result, err := PostFormRequest(server.URL, tt.data)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.response, string(result))
			}
		})
	}
}

func TestGetRequest(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		statusCode int
		response   string
		expectErr  bool
	}{
		{
			name:       "successful get",
			token:      "valid-token",
			statusCode: http.StatusOK,
			response:   "success response",
			expectErr:  false,
		},
		{
			name:       "unauthorized",
			token:      "invalid-token",
			statusCode: http.StatusUnauthorized,
			response:   "unauthorized",
			expectErr:  false,
		},
		{
			name:       "empty token",
			token:      "",
			statusCode: http.StatusUnauthorized,
			response:   "unauthorized",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "Bearer "+tt.token, r.Header.Get("Authorization"))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			result, err := GetRequest(server.URL, tt.token)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.response, string(result))
			}
		})
	}
}
