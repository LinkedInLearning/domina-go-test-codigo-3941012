package database

import (
	"context"

	"pokemon-battle/internal/models"
)

type battleService struct {
	// BattleCRUDService is a generic CRUD service implemented by the service
	BattleCRUDService

	// srv is the service with the actual database connection
	srv Service
}

func NewBattleService() *battleService {
	return &battleService{
		srv: New(),
	}
}

// Create inserts a new battle into the database
func (s *battleService) Create(ctx context.Context, battle *models.Battle) error {
	db := s.srv.MustDB()

	if err := battle.Validate(); err != nil {
		return err
	}

	query := "INSERT INTO battles (pokemon1_id, pokemon2_id, winner_id) VALUES ($1, $2, $3) RETURNING id"

	return db.QueryRowContext(ctx, query, battle.Pokemon1ID, battle.Pokemon2ID, battle.WinnerID).Scan(&battle.ID)
}

// DeleteBattle deletes a battle from the database
func (s *battleService) Delete(ctx context.Context, id int) error {
	db := s.srv.MustDB()

	query := "DELETE FROM battles WHERE id=$1"
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// GetAll retrieves all battles from the database
func (s *battleService) GetAll(ctx context.Context) ([]models.Battle, error) {
	db := s.srv.MustDB()

	query := "SELECT id, pokemon1_id, pokemon2_id, winner_id FROM battles"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []models.Battle
	for rows.Next() {
		var battle models.Battle
		if err := rows.Scan(&battle.ID, &battle.Pokemon1ID, &battle.Pokemon2ID, &battle.WinnerID); err != nil {
			return nil, err
		}
		battles = append(battles, battle)
	}

	return battles, nil
}

// GetByID retrieves a battle from the database by its ID
func (s *battleService) GetByID(ctx context.Context, id int) (models.Battle, error) {
	db := s.srv.MustDB()

	query := "SELECT id, pokemon1_id, pokemon2_id, winner_id FROM battles WHERE id=$1"
	row := db.QueryRowContext(ctx, query, id)

	var battle models.Battle
	if err := row.Scan(&battle.ID, &battle.Pokemon1ID, &battle.Pokemon2ID, &battle.WinnerID); err != nil {
		return models.Battle{}, err
	}
	return battle, nil
}

// Update updates an existing battle in the database
func (s *battleService) Update(ctx context.Context, battle models.Battle) error {
	db := s.srv.MustDB()

	if err := battle.Validate(); err != nil {
		return err
	}

	query := "UPDATE battles SET pokemon1_id=$1, pokemon2_id=$2, winner_id=$3 WHERE id=$4"
	_, err := db.ExecContext(ctx, query, battle.Pokemon1ID, battle.Pokemon2ID, battle.WinnerID, battle.ID)
	return err
}
