package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPokemonService is used for testing the pokemon routes
// Implementes the PokemonService interface using testify mock.
type mockPokemonService struct {
	mock.Mock
}

func (m *mockPokemonService) Create(ctx context.Context, pokemon *models.Pokemon) error {
	args := m.Called(ctx, pokemon)
	return args.Error(0)
}

func (m *mockPokemonService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPokemonService) GetAll(ctx context.Context) ([]models.Pokemon, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Pokemon), args.Error(1)
}

func (m *mockPokemonService) GetByID(ctx context.Context, id int) (models.Pokemon, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Pokemon), args.Error(1)
}

func (m *mockPokemonService) Update(ctx context.Context, pokemon models.Pokemon) error {
	args := m.Called(ctx, pokemon)
	return args.Error(0)
}

func TestGetAllPokemons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()

		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("GetAll", mock.Anything).Return([]models.Pokemon{
			{ID: 1, Name: "Pikachu", Type: "Electric", HP: 10, Attack: 10, Defense: 10},
			{ID: 2, Name: "Charmander", Type: "Fire", HP: 10, Attack: 10, Defense: 10},
		}, nil)

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/pokemons", nil)
		require.NoError(t, err)
		// Perform the request
		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var pokemons []models.Pokemon
		err = json.Unmarshal(body, &pokemons)
		require.NoError(t, err)
		require.Equal(t, len(pokemons), 2)

		require.Equal(t, pokemons[0].ID, 1)
		require.Equal(t, pokemons[0].Name, "Pikachu")
		require.Equal(t, pokemons[0].Type, "Electric")
		require.Equal(t, pokemons[0].HP, 10)
		require.Equal(t, pokemons[0].Attack, 10)
		require.Equal(t, pokemons[0].Defense, 10)

		require.Equal(t, pokemons[1].ID, 2)
		require.Equal(t, pokemons[1].Name, "Charmander")
		require.Equal(t, pokemons[1].Type, "Fire")
		require.Equal(t, pokemons[1].HP, 10)
		require.Equal(t, pokemons[1].Attack, 10)
		require.Equal(t, pokemons[1].Defense, 10)
	})

	t.Run("error", func(t *testing.T) {
		s := New()

		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("GetAll", mock.Anything).Return([]models.Pokemon{}, errors.New("mock error"))

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/pokemons", nil)
		require.NoError(t, err)
		// Perform the request
		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestCreatePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Create", mock.Anything, mock.Anything).Return(nil)

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Post("/", pokemonServer.CreatePokemon)

		pokemon := models.Pokemon{
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}
		body, err := json.Marshal(pokemon)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/pokemons", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusCreated)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Create", mock.Anything, mock.Anything).Return(errors.New("mock error"))

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Post("/", pokemonServer.CreatePokemon)

		pokemon := models.Pokemon{
			Name: "Bulbasaur",
			Type: "Grass",
		}
		body, err := json.Marshal(pokemon)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/pokemons", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestGetPokemonByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("GetByID", mock.Anything, mock.Anything).Return(models.Pokemon{
			ID:      1,
			Name:    "Pikachu",
			Type:    "Electric",
			HP:      10,
			Attack:  10,
			Defense: 10,
		}, nil)

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)

		req, err := http.NewRequest("GET", "/pokemons/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("GetByID", mock.Anything, mock.Anything).Return(models.Pokemon{}, errors.New("mock error"))

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)

		req, err := http.NewRequest("GET", "/pokemons/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestUpdatePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Update", mock.Anything, mock.Anything).Return(nil)

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)

		pokemon := models.Pokemon{
			ID:      1,
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}
		body, _ := json.Marshal(pokemon)

		req, err := http.NewRequest("PUT", "/pokemons/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Update", mock.Anything, mock.Anything).Return(errors.New("mock error"))

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)

		pokemon := models.Pokemon{
			ID:   1,
			Name: "Bulbasaur",
			Type: "Grass",
		}
		body, _ := json.Marshal(pokemon)

		req, err := http.NewRequest("PUT", "/pokemons/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestDeletePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Delete", mock.Anything, mock.Anything).Return(nil)

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

		req, err := http.NewRequest("DELETE", "/pokemons/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		mockPSrv := &mockPokemonService{}
		mockPSrv.On("Delete", mock.Anything, mock.Anything).Return(errors.New("mock error"))

		pokemonServer := pokemonServer{srv: mockPSrv}
		pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

		req, err := http.NewRequest("DELETE", "/pokemons/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}
