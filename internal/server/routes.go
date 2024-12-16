package server

import (
	"pokemon-battle/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	// init the pokemon routes from a pokemon service
	pokemonServer := pokemonServer{srv: database.NewPokemonService()}

	pokemonRoutes := s.App.Group("/pokemons")
	pokemonRoutes.Post("/", pokemonServer.CreatePokemon)
	pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)
	pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)
	pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)
	pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

	// init the battle routes from a battle service
	battleServer := battleServer{srv: database.NewBattleService()}

	battleRoutes := s.App.Group("/battles")
	battleRoutes.Post("/", battleServer.CreateBattle)
	battleRoutes.Get("/", battleServer.GetAllBattles)
	battleRoutes.Get("/:id", battleServer.GetBattleByID)
	battleRoutes.Put("/:id", battleServer.UpdateBattle)
	battleRoutes.Delete("/:id", battleServer.DeleteBattle)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Welcome to Pokemon Battle!",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
