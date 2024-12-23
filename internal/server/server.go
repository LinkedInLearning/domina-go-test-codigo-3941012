package server

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/business"
	"pokemon-battle/internal/database"
)

type FiberServer struct {
	*fiber.App

	db        database.Service
	diceSides int
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "pokemon-battle",
			AppName:      "pokemon-battle",
		}),

		db:        database.New(),
		diceSides: initalizeDiceSides(),
	}

	return server
}

func initalizeDiceSides() int {
	sides, err := strconv.Atoi(os.Getenv("POKEMON_BATTLE_DICE_SIDES"))
	if err != nil {
		return business.DefaultDiceSides
	}
	return sides
}
