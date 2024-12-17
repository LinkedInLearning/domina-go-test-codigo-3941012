package business

import (
	"math/rand"
)

const DefaultDiceSides = 6

// savageDice tira un dado salvaje, que explota cuando sale el máximo valor.
// El dado explota repitiendo la tirada hasta que deje de explotar,
// o hasta que se hayan hecho 50 tiradas, que es el máximo de tiradas permitidas.
func savageDice(diceSides int) int {
	sum := 0
	maxRolls := 50

	// Initializar la tirada con el máximo valor
	// para que entrar a tirar.
	roll := diceSides

	// Un dado explota cuando sale el máximo valor,
	// repitiendo la tirada hasta que deje de explotar.
	for roll == diceSides && maxRolls > 0 {
		maxRolls--
		roll = rand.Intn(diceSides) + 1
		sum += roll
	}

	return sum
}
