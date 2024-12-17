package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/business"
	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

// battleServer is used to handle the battle routes.
// It receives a database.BattleCRUDService and uses it to handle the routes.
type battleServer struct {
	srv        database.BattleCRUDService
	pokemonSrv database.PokemonCRUDService
	diceSides  int
}

type battleRequest struct {
	Pokemon1ID int `json:"pokemon1_id"`
	Pokemon2ID int `json:"pokemon2_id"`
}

func (s *battleServer) CreateBattle(c *fiber.Ctx) error {
	ctx := context.Background()
	var req battleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// retrieve the pokemons from the database
	pokemon1, err := s.pokemonSrv.GetByID(ctx, req.Pokemon1ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	pokemon2, err := s.pokemonSrv.GetByID(ctx, req.Pokemon2ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	battle := business.Fight(s.diceSides, pokemon1, pokemon2)

	err = s.srv.Create(ctx, &battle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(battle)
}

func (s *battleServer) GetAllBattles(c *fiber.Ctx) error {
	ctx := context.Background()
	battles, err := s.srv.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battles)
}

func (s *battleServer) GetBattleByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	battle, err := s.srv.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battle)
}

func (s *battleServer) UpdateBattle(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var battle models.Battle
	if err := c.BodyParser(&battle); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	battle.ID = id

	err = s.srv.Update(ctx, battle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battle)
}

func (s *battleServer) DeleteBattle(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = s.srv.Delete(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
