package business

import (
	"pokemon-battle/internal/models"
)

const initiativeDiceSides = 100

func Fight(diceSides int, pokemon1 models.Pokemon, pokemon2 models.Pokemon) (models.Battle, error) {
	// Create a battle record
	battle := models.Battle{
		Pokemon1ID: pokemon1.ID,
		Pokemon2ID: pokemon2.ID,
	}

	// Battle continues until one Pokemon's HP reaches 0
	turns := 1
	for {
		// Decide who starts (1-100 roll)
		var startRoll1, startRoll2 int
		for startRoll1 == startRoll2 {
			startRoll1 = rollDice(initiativeDiceSides)
			startRoll2 = rollDice(initiativeDiceSides)
		}

		attacker, defender := &pokemon1, &pokemon2
		if startRoll2 > startRoll1 {
			attacker, defender = &pokemon2, &pokemon1
		}

		attack(diceSides, attacker, defender)

		// If defender is still alive, they get to attack
		if defender.HP > 0 {
			attack(diceSides, defender, attacker)
		}

		// Determine winner, if one of them is without HP
		if attacker.HP <= 0 {
			battle.WinnerID = defender.ID
			break
		} else if defender.HP <= 0 {
			battle.WinnerID = attacker.ID
			break
		}

		turns++
	}

	battle.Turns = turns

	return battle, nil
}

func attack(diceSides int, attacker *models.Pokemon, defender *models.Pokemon) {
	// Calculate attack value (base attack + dice roll)
	attackRoll := rollDice(diceSides)
	totalAttack := attacker.Attack + attackRoll

	// Calculate defense value (base defense + dice roll)
	defenseRoll := rollDice(diceSides)
	totalDefense := defender.Defense + defenseRoll

	// If attack beats defense, reduce defender's HP
	if totalAttack > totalDefense {
		damage := totalAttack - totalDefense
		defender.HP -= damage
	}
}
