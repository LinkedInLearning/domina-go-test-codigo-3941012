package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func (s *FiberServer) CreateBattle(c *fiber.Ctx) error {
	srv := database.NewBattleService()

	ctx := context.Background()
	var battle models.Battle
	if err := c.BodyParser(&battle); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := srv.Create(ctx, &battle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(battle)
}

func (s *FiberServer) GetAllBattles(c *fiber.Ctx) error {
	srv := database.NewBattleService()

	ctx := context.Background()
	battles, err := srv.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battles)
}

func (s *FiberServer) GetBattleByID(c *fiber.Ctx) error {
	srv := database.NewBattleService()

	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	battle, err := srv.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battle)
}

func (s *FiberServer) UpdateBattle(c *fiber.Ctx) error {
	srv := database.NewBattleService()

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
	err = srv.Update(ctx, battle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(battle)
}

func (s *FiberServer) DeleteBattle(c *fiber.Ctx) error {
	srv := database.NewBattleService()

	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	err = srv.Delete(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
