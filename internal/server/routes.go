package server

import (
	"pokemon-battle/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes(pokemonSrv database.PokemonCRUDService, battleSrv database.BattleCRUDService) {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Apply Basic Auth middleware, only ash, misty and brock are allowed
	// to access all the routes
	s.App.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"ash":   "ketchum",
			"misty": "waters",
			"brock": "brock",
		},
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	// init the pokemon routes from a pokemon service
	pokemonServer := pokemonServer{srv: pokemonSrv}

	pokemonRoutes := s.App.Group("/pokemons")
	pokemonRoutes.Post("/", pokemonServer.CreatePokemon)
	pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)
	pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)
	pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)
	pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

	// init the battle routes from a battle service
	battleServer := battleServer{srv: battleSrv, pokemonSrv: pokemonSrv, diceSides: s.diceSides}

	battleRoutes := s.App.Group("/battles")
	battleRoutes.Post("/", battleServer.CreateBattle)
	battleRoutes.Get("/", battleServer.GetAllBattles)
	battleRoutes.Get("/:id", battleServer.GetBattleByID)
	battleRoutes.Put("/:id", battleServer.UpdateBattle)
	battleRoutes.Delete("/:id", battleServer.DeleteBattle)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	// get the username and password from the context
	username := c.Locals("username").(string)

	resp := fiber.Map{
		"message":  "Welcome to Pokemon Battle!",
		"username": username,
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
