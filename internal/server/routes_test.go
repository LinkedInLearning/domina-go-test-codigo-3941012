package server

import (
	"io"
	"net/http"
	"pokemon-battle/internal/models"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	s := New()

	mockPSrv := &mockPokemonService{}
	mockPSrv.On("GetAll", mock.Anything).Return([]models.Pokemon{}, nil)
	mockPSrv.On("GetByID", mock.Anything, mock.Anything).Return(models.Pokemon{}, nil)
	mockPSrv.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockPSrv.On("Delete", mock.Anything, mock.Anything).Return(nil)

	mockBSrv := &mockBattleService{}
	mockBSrv.On("GetAll", mock.Anything).Return([]models.Battle{}, nil)
	mockBSrv.On("GetByID", mock.Anything, mock.Anything).Return(models.Battle{}, nil)
	mockBSrv.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockBSrv.On("Delete", mock.Anything, mock.Anything).Return(nil)

	s.RegisterFiberRoutes(mockPSrv, mockBSrv)

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
