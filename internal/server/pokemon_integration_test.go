package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func TestPokemon_IT(t *testing.T) {
	s := New()

	databaseSrv := MustNewWithDatabase()
	pokemonSrv := database.NewPokemonService(databaseSrv)

	// the battle service is mocked because it's not needed for this test
	s.RegisterFiberRoutes(pokemonSrv, &mockBattleService{hasError: true})

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
			if err != nil {
				t.Fatalf("error making request to server. Err: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected status Created; got %v", resp.Status)
			}

			var pokemonResponse models.Pokemon
			err = json.NewDecoder(resp.Body).Decode(&pokemonResponse)
			if err != nil {
				t.Fatalf("error decoding response. Err: %v", err)
			}

			if pokemonResponse.ID == 0 {
				t.Errorf("expected pokemon ID to be greater than 0; got %v", pokemonResponse.ID)
			}
			if pokemonResponse.Name != pokemonReq.Name ||
				pokemonResponse.Type != pokemonReq.Type ||
				pokemonResponse.HP != pokemonReq.HP ||
				pokemonResponse.Attack != pokemonReq.Attack ||
				pokemonResponse.Defense != pokemonReq.Defense {

				t.Errorf("expected pokemon to be %v; got %v", pokemonReq, pokemonResponse)
			}
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
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var pokemonResponse models.Pokemon
		err = json.NewDecoder(resp.Body).Decode(&pokemonResponse)
		if err != nil {
			t.Fatalf("error decoding response. Err: %v", err)
		}

		if pokemonResponse.ID != p.ID {
			t.Errorf("expected pokemon ID to be %v; got %v", p.ID, pokemonResponse.ID)
		}
	})

	t.Run("get-all", func(t *testing.T) {
		initialPokemons, err := pokemonSrv.GetAll(context.Background())
		if err != nil {
			t.Fatalf("error getting pokemons. Err: %v", err)
		}

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
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var pokemons []models.Pokemon
		err = json.NewDecoder(resp.Body).Decode(&pokemons)
		if err != nil {
			t.Fatalf("error decoding response. Err: %v", err)
		}

		if len(pokemons) != len(initialPokemons)+10 {
			t.Errorf("expected %v pokemons; got %v", len(initialPokemons)+10, len(pokemons))
		}
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
		if err != nil {
			t.Fatalf("error marshalling pokemon. Err: %v", err)
		}

		req := createAuthenticatedRequest(t, "PUT", "/pokemons/"+strconv.Itoa(p.ID), body)

		resp, err := s.App.Test(req, -1) // disable timeout
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		// get the pokemon from the db
		pokemon, err := pokemonSrv.GetByID(context.Background(), p.ID)
		if err != nil {
			t.Fatalf("error getting pokemon. Err: %v", err)
		}

		if pokemon.HP != p.HP || pokemon.Attack != p.Attack || pokemon.Defense != p.Defense {
			t.Errorf("expected HP to be %v, Attack to be %v, and Defense to be %v; got %v, %v, and %v", p.HP, p.Attack, p.Defense, pokemon.HP, pokemon.Attack, pokemon.Defense)
		}
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
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status NoContent; got %v", resp.Status)
		}

		// get the pokemon from the db
		pokemon, err := pokemonSrv.GetByID(context.Background(), p.ID)
		if err == nil {
			t.Fatalf("expected error getting pokemon. Err: %v", err)
		}

		if pokemon.ID != 0 {
			t.Errorf("expected pokemon ID to be 0; got %v", pokemon.ID)
		}
	})
}

// createPokemonUsingDB is a helper function to create a pokemon using the database service
func createPokemonUsingDB(t *testing.T, srv database.PokemonCRUDService, p *models.Pokemon) {
	t.Helper()

	err := srv.Create(context.Background(), p)
	if err != nil {
		t.Fatalf("error inserting pokemon. Err: %v", err)
	}
}
