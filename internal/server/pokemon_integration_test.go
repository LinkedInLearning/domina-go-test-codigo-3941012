package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/gofiber/fiber/v2"
)

func TestPokemon_IT(t *testing.T) {
	app := fiber.New()
	s := &FiberServer{App: app}

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

			req, err := http.NewRequest("POST", "/pokemons", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
			if err != nil {
				t.Fatalf("error creating request. Err: %v", err)
			}

			resp, err := app.Test(req, -1) // disable timeout
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

		// use the db layer to insert a pokemon
		err := pokemonSrv.Create(context.Background(), &p)
		if err != nil {
			t.Fatalf("error inserting pokemon. Err: %v", err)
		}

		req, err := http.NewRequest("GET", "/pokemons/"+strconv.Itoa(p.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
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

			// use the db layer to insert a battle
			err := pokemonSrv.Create(context.Background(), &p)
			if err != nil {
				t.Fatalf("error inserting pokemon. Err: %v", err)
			}
		}

		req, err := http.NewRequest("GET", "/pokemons", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
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

		// use the db layer to insert a pokemon
		err := pokemonSrv.Create(context.Background(), &p)
		if err != nil {
			t.Fatalf("error inserting pokemon. Err: %v", err)
		}

		// update the pokemon
		p.HP = 50
		p.Attack = 51
		p.Defense = 51

		body, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("error marshalling pokemon. Err: %v", err)
		}

		req, err := http.NewRequest("PUT", "/pokemons/"+strconv.Itoa(p.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
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

		// use the db layer to insert a pokemon
		err := pokemonSrv.Create(context.Background(), &p)
		if err != nil {
			t.Fatalf("error inserting pokemon. Err: %v", err)
		}

		req, err := http.NewRequest("DELETE", "/pokemons/"+strconv.Itoa(p.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
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
