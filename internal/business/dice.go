package business

import (
	"math/rand"
)

const DefaultDiceSides = 6

func rollDice(diceSides int) int {
	return rand.Intn(diceSides) + 1
}
