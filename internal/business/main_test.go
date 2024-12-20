package business_test

import (
	"os"
	"testing"

	"pokemon-battle/internal/models"
)

var (
	// weakPokemon es un pokemon con poca HP, Attack y Defense.
	// Usado en los tests para verificar el pokemon perdedor.
	weakPokemon models.Pokemon

	// strongPokemon es un pokemon con mucha HP, Attack y Defense.
	// Usado en los tests para verificar el pokemon ganador.
	strongPokemon models.Pokemon
)

func TestMain(m *testing.M) {

	// inicializar los pokemons usados en los tests una vez.
	weakPokemon = models.Pokemon{
		ID:      1,
		Name:    "Pikachu",
		Type:    "Electric",
		HP:      1,
		Attack:  1,
		Defense: 1,
	}

	strongPokemon = models.Pokemon{
		ID:      2,
		Name:    "Charizard",
		Type:    "Fire",
		HP:      100,
		Attack:  55,
		Defense: 40,
	}

	os.Exit(m.Run())
}
