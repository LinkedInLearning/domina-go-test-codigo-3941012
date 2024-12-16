package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

// battleServer is used to handle the battle routes.
// It receives a database.BattleCRUDService and uses it to handle the routes.
type battleServer struct {
	srv database.BattleCRUDService
}

func (s *battleServer) CreateBattle(c *fiber.Ctx) error {
	ctx := context.Background()
	var battle models.Battle
	if err := c.BodyParser(&battle); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := s.srv.Create(ctx, &battle)
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
