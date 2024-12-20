package business

import (
	"math/rand"
)

const DefaultDiceSides = 6

// Dice es una interfaz que representa un dado.
type Dice interface {
	Roll() int
	Result() int
}

// BaseDice es una implementación de Dice que representa un dado base.
type BaseDice struct {
	Sides  int
	result int
}

// Result devuelve el resultado total de la tirada.
func (d *BaseDice) Result() int {
	return d.result
}

// SavageDice es una implementación de Dice que representa un dado salvaje,
// esto es, un dado que puede explotar hasta 50 veces.
type SavageDice struct {
	BaseDice
	maxExplosions int
	Explosions    int
}

func (d *BaseDice) Roll() int {
	d.result = rand.Intn(d.Sides) + 1
	return d.result
}

func (d *SavageDice) Roll() int {
	sum := 0
	maxRolls := d.maxExplosions
	if maxRolls == 0 {
		// Si no se ha definido el máximo de explosiones,
		// se asume que el dado no explota.
		maxRolls = 1
	}

	// Initializar la tirada con el máximo valor
	// para que entrar a tirar.
	roll := d.Sides

	// Un dado explota cuando sale el máximo valor,
	// repitiendo la tirada hasta que deje de explotar.
	explosions := 0
	for roll == d.Sides && maxRolls > 0 {
		maxRolls--
		roll = rand.Intn(d.Sides) + 1
		sum += roll
		if roll == d.Sides {
			explosions++
		}
	}

	d.result = sum
	d.Explosions = explosions

	return sum
}
