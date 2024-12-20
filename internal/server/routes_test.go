package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	s := New()

	s.RegisterFiberRoutes(&mockPokemonService{hasError: false}, &mockBattleService{hasError: false})

	t.Run("get/", func(t *testing.T) {
		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // ash:ketchum in base64
		require.NoError(t, err)

		// Perform the request
		resp, err := s.App.Test(req)
		require.NoError(t, err)
		// Your test assertions...
		require.Equal(t, resp.StatusCode, http.StatusOK)

		expected := "{\"message\":\"Welcome to Pokemon Battle!\",\"username\":\"ash\"}"
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, expected, string(body))
	})
}
