package server

import (
	"io"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	s := New()

	s.RegisterFiberRoutes(&mockPokemonService{hasError: false}, &mockBattleService{hasError: false})

	t.Run("get/", func(t *testing.T) {
		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // ash:ketchum in base64
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		// Perform the request
		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		// Your test assertions...
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		expected := "{\"message\":\"Welcome to Pokemon Battle!\",\"username\":\"ash\"}"
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading response body. Err: %v", err)
		}
		if expected != string(body) {
			t.Errorf("expected response body to be %v; got %v", expected, string(body))
		}
	})
}
