package models

import "errors"

type Pokemon struct {
	ID      int    `json:"id"`      // Identificador único del Pokémon
	Name    string `json:"name"`    // Nombre del Pokémon
	Type    string `json:"type"`    // Tipo del Pokémon (e.g., "Fuego", "Agua")
	HP      int    `json:"hp"`      // Puntos de salud
	Attack  int    `json:"attack"`  // Nivel de ataque
	Defense int    `json:"defense"` // Nivel de defensa
}

func (p *Pokemon) Validate() error {
	if p.Name == "" {
		return errors.New("pokemon name cannot be empty")
	}
	if p.Type == "" {
		return errors.New("pokemon type cannot be empty")
	}
	if p.HP <= 0 {
		return errors.New("pokemon HP must be greater than 0")
	}
	if p.Attack < 0 {
		return errors.New("pokemon attack cannot be negative")
	}
	if p.Defense < 0 {
		return errors.New("pokemon defense cannot be negative")
	}
	return nil
}

type Battle struct {
	ID         int `json:"id"`          // Identificador único de la batalla
	Pokemon1ID int `json:"pokemon1_id"` // ID del primer Pokémon participante
	Pokemon2ID int `json:"pokemon2_id"` // ID del segundo Pokémon participante
	WinnerID   int `json:"winner_id"`   // ID del Pokémon ganador
}

func (b *Battle) Validate() error {
	if b.Pokemon1ID <= 0 || b.Pokemon2ID <= 0 {
		return errors.New("invalid pokemon IDs")
	}

	if b.Pokemon1ID == b.Pokemon2ID {
		return errors.New("pokemon cannot battle itself")
	}

	if b.WinnerID != b.Pokemon1ID && b.WinnerID != b.Pokemon2ID {
		return errors.New("winner must be one of the battling pokemon")
	}
	return nil
}
