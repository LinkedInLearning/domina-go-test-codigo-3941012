package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPokemon_IT(t *testing.T) {
	s := New()

	databaseSrv := MustNewWithDatabase(t)
	pokemonSrv := database.NewPokemonService(databaseSrv)

	// the battle service is mocked because it's not needed for this test
	mockBSrv := &mockBattleService{}
	mockBSrv.On("GetAll", mock.Anything).Return([]models.Battle{}, nil)
	mockBSrv.On("GetByID", mock.Anything, mock.Anything).Return(models.Battle{}, nil)
	mockBSrv.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockBSrv.On("Delete", mock.Anything, mock.Anything).Return(nil)

	s.RegisterFiberRoutes(pokemonSrv, mockBSrv)

	t.Run("create", func(t *testing.T) {
		t.Run("post-ok", func(t *testing.T) {
			pokemonReq := pokemonRequest{
				Name:    "Bulbasaur",
				Type:    "Grass",
				HP:      45,
				Attack:  49,
				Defense: 49,
			}
			body, _ := json.Marshal(pokemonReq)

			req := createAuthenticatedRequest(t, "POST", "/pokemons", body)

			resp, err := s.App.Test(req, -1) // disable timeout
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, resp.StatusCode, http.StatusCreated)

			var pokemonResponse models.Pokemon
			err = json.NewDecoder(resp.Body).Decode(&pokemonResponse)
			require.NoError(t, err)

			require.Greater(t, pokemonResponse.ID, 0)
			require.Equal(t, pokemonResponse.Name, pokemonReq.Name)
			require.Equal(t, pokemonResponse.Type, pokemonReq.Type)
			require.Equal(t, pokemonResponse.HP, pokemonReq.HP)
			require.Equal(t, pokemonResponse.Attack, pokemonReq.Attack)
			require.Equal(t, pokemonResponse.Defense, pokemonReq.Defense)
		})
	})

	t.Run("get-by-id", func(t *testing.T) {
		p := models.Pokemon{
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}

		createPokemonUsingDB(t, pokemonSrv, &p)

		req := createAuthenticatedRequest(t, "GET", "/pokemons/"+strconv.Itoa(p.ID), nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		var pokemonResponse models.Pokemon
		err = json.NewDecoder(resp.Body).Decode(&pokemonResponse)
		require.NoError(t, err)

		require.Equal(t, pokemonResponse.ID, p.ID)
	})

	t.Run("get-all", func(t *testing.T) {
		initialPokemons, err := pokemonSrv.GetAll(context.Background())
		require.NoError(t, err)

		// insert 10 pokemons
		for i := 0; i < 10; i++ {
			p := models.Pokemon{
				Name:    "Bulbasaur",
				Type:    "Grass",
				HP:      45,
				Attack:  49,
				Defense: 49,
			}

			createPokemonUsingDB(t, pokemonSrv, &p)
		}

		req := createAuthenticatedRequest(t, "GET", "/pokemons", nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		var pokemons []models.Pokemon
		err = json.NewDecoder(resp.Body).Decode(&pokemons)
		require.NoError(t, err)

		require.Equal(t, len(pokemons), len(initialPokemons)+10)
	})

	t.Run("update", func(t *testing.T) {
		p := models.Pokemon{
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}

		createPokemonUsingDB(t, pokemonSrv, &p)

		// update the pokemon
		p.HP = 50
		p.Attack = 51
		p.Defense = 51

		body, err := json.Marshal(p)
		require.NoError(t, err)

		req := createAuthenticatedRequest(t, "PUT", "/pokemons/"+strconv.Itoa(p.ID), body)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		// get the pokemon from the db
		pokemon, err := pokemonSrv.GetByID(context.Background(), p.ID)
		require.NoError(t, err)

		require.Equal(t, pokemon.HP, p.HP)
		require.Equal(t, pokemon.Attack, p.Attack)
		require.Equal(t, pokemon.Defense, p.Defense)
	})

	t.Run("delete", func(t *testing.T) {
		p := models.Pokemon{
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}

		createPokemonUsingDB(t, pokemonSrv, &p)

		req := createAuthenticatedRequest(t, "DELETE", "/pokemons/"+strconv.Itoa(p.ID), nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusNoContent)

		// get the pokemon from the db
		pokemon, err := pokemonSrv.GetByID(context.Background(), p.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.Equal(t, pokemon.ID, 0)
	})
}

// createPokemonUsingDB is a helper function to create a pokemon using the database service
func createPokemonUsingDB(t *testing.T, srv database.PokemonCRUDService, p *models.Pokemon) {
	t.Helper()

	err := srv.Create(context.Background(), p)
	require.NoError(t, err)
}
