package database

import (
	"context"

	"pokemon-battle/internal/models"
)

type pokemonService struct {
	// PokemonCRUDService is a generic CRUD service implemented by the service
	PokemonCRUDService

	// srv is the service with the actual database connection
	srv Service
}

func NewPokemonService() *pokemonService {
	return &pokemonService{
		srv: New(),
	}
}

// Create inserts a new pokemon into the database
func (s *pokemonService) Create(ctx context.Context, pokemon *models.Pokemon) error {
	db := s.srv.MustDB()

	if err := pokemon.Validate(); err != nil {
		return err
	}

	query := "INSERT INTO pokemons (name, type, hp, attack, defense) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	return db.QueryRowContext(ctx, query, pokemon.Name, pokemon.Type, pokemon.HP, pokemon.Attack, pokemon.Defense).Scan(&pokemon.ID)
}

// Delete deletes a pokemon from the database
func (s *pokemonService) Delete(ctx context.Context, id int) error {
	db := s.srv.MustDB()

	query := "DELETE FROM pokemons WHERE id=$1"
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// GetAll retrieves all pokemons from the database
func (s *pokemonService) GetAll(ctx context.Context) ([]models.Pokemon, error) {
	db := s.srv.MustDB()

	query := "SELECT id, name, type, hp, attack, defense FROM pokemons ORDER BY id"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokemons []models.Pokemon
	for rows.Next() {
		var pokemon models.Pokemon
		if err := rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Type, &pokemon.HP, &pokemon.Attack, &pokemon.Defense); err != nil {
			return nil, err
		}
		pokemons = append(pokemons, pokemon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pokemons, nil
}

// GetByID retrieves a pokemon from the database by its ID
func (s *pokemonService) GetByID(ctx context.Context, id int) (models.Pokemon, error) {
	db := s.srv.MustDB()

	query := "SELECT id, name, type, hp, attack, defense FROM pokemons WHERE id=$1"
	row := db.QueryRowContext(ctx, query, id)

	var pokemon models.Pokemon
	if err := row.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Type, &pokemon.HP, &pokemon.Attack, &pokemon.Defense); err != nil {
		return models.Pokemon{}, err
	}
	return pokemon, nil
}

// Update updates an existing pokemon in the database
func (s *pokemonService) Update(ctx context.Context, pokemon models.Pokemon) error {
	db := s.srv.MustDB()

	if err := pokemon.Validate(); err != nil {
		return err
	}

	query := "UPDATE pokemons SET name=$1, type=$2, hp=$3, attack=$4, defense=$5 WHERE id=$6"
	_, err := db.ExecContext(ctx, query, pokemon.Name, pokemon.Type, pokemon.HP, pokemon.Attack, pokemon.Defense, pokemon.ID)
	return err
}
