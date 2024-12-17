package server

import (
	"github.com/gofiber/fiber/v2"

	"pokemon-battle/internal/database"
)

type FiberServer struct {
	*fiber.App

	db        database.Service
	diceSides int
}

func New(diceSides int) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "pokemon-battle",
			AppName:      "pokemon-battle",
		}),

		db:        database.New(),
		diceSides: diceSides,
	}

	return server
}
