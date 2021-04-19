package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_Wrap(t *testing.T) {
	// Given
	s := NewServer()

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	s.Wrap(http.MethodGet, "/test", func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	// When
	resp, err := http.Get(fmt.Sprintf("%s/test", ts.URL))
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServer_Wrap_HandleError(t *testing.T) {
	tt := []struct {
		name               string
		err                error
		expectedCode       string
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			name:               "bad request",
			err:                NewError("error", http.StatusBadRequest),
			expectedCode:       "bad_request",
			expectedMessage:    "error",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "not found",
			err:                NewError("error", http.StatusNotFound),
			expectedCode:       "not_found",
			expectedMessage:    "error",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "an error but not my error",
			err:                errors.New("an error but not my error"),
			expectedCode:       "internal_server_error",
			expectedMessage:    "an error but not my error",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			s := NewServer()

			ts := httptest.NewServer(s.Router)
			defer ts.Close()

			s.Wrap(http.MethodGet, "/test", func(w http.ResponseWriter, r *http.Request) error {
				return tc.err
			})

			// When
			resp, err := http.Get(fmt.Sprintf("%s/test", ts.URL))
			if err != nil {
				t.Fatal(err)
			}

			var r struct {
				StatusCode int    `json:"status"`
				Code       string `json:"code"`
				Message    string `json:"message"`
			}

			_ = json.NewDecoder(resp.Body).Decode(&r)

			// Then
			require.Equal(t, tc.expectedCode, r.Code)
			require.Equal(t, tc.expectedStatusCode, r.StatusCode)
			require.Equal(t, tc.expectedMessage, r.Message)
		})
	}
}
