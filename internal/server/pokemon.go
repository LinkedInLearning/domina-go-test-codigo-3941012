package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func (s *FiberServer) CreatePokemon(c *fiber.Ctx) error {
	srv := database.NewPokemonService()

	ctx := context.Background()
	var pokemon models.Pokemon
	if err := c.BodyParser(&pokemon); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := srv.Create(ctx, &pokemon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(pokemon)
}

func (s *FiberServer) GetAllPokemons(c *fiber.Ctx) error {
	srv := database.NewPokemonService()

	ctx := context.Background()
	pokemons, err := srv.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemons)
}

func (s *FiberServer) GetPokemonByID(c *fiber.Ctx) error {
	srv := database.NewPokemonService()

	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pokemon, err := srv.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemon)
}

func (s *FiberServer) UpdatePokemon(c *fiber.Ctx) error {
	srv := database.NewPokemonService()

	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var pokemon models.Pokemon
	if err := c.BodyParser(&pokemon); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	pokemon.ID = id
	err = srv.Update(ctx, pokemon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemon)
}

func (s *FiberServer) DeletePokemon(c *fiber.Ctx) error {
	srv := database.NewPokemonService()

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
