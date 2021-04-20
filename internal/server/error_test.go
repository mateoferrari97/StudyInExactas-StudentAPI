package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	// Given
	e := NewError("error", http.StatusInternalServerError)

	// When
	m := e.Error()

	// Then
	require.Equal(t, "500 internal_server_error: error", m)
}

func TestError_EmptyValues(t *testing.T) {
	tt := []struct {
		name          string
		err           error
		expectedError string
	}{
		{
			name:          "0 status code",
			err:           NewError("error", 0),
			expectedError: "error",
		},
		{
			name:          "empty message",
			err:           NewError("", http.StatusInternalServerError),
			expectedError: "internal_server_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			e := tc.err

			// When
			m := e.Error()

			// Then
			require.Equal(t, tc.expectedError, m)
		})
	}
}
