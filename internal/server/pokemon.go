package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

// pokemonServer is used to handle the pokemon routes.
// It receives a database.PokemonCRUDService and uses it to handle the routes.
type pokemonServer struct {
	srv database.PokemonCRUDService
}

func (s *pokemonServer) CreatePokemon(c *fiber.Ctx) error {
	ctx := context.Background()
	var pokemon models.Pokemon
	if err := c.BodyParser(&pokemon); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := s.srv.Create(ctx, &pokemon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(pokemon)
}

func (s *pokemonServer) GetAllPokemons(c *fiber.Ctx) error {
	ctx := context.Background()
	pokemons, err := s.srv.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemons)
}

func (s *pokemonServer) GetPokemonByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pokemon, err := s.srv.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemon)
}

func (s *pokemonServer) UpdatePokemon(c *fiber.Ctx) error {
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

	err = s.srv.Update(ctx, pokemon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pokemon)
}

func (s *pokemonServer) DeletePokemon(c *fiber.Ctx) error {
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
