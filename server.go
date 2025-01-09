package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"survey-service/config"
	"survey-service/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	ctx := context.Background()
	routes.SetRoutes(ctx, app)

	go func() {
		log.Println(fmt.Sprintf("Starting server on port %s", config.PprofPort))
		log.Println(http.ListenAndServe(":6060", nil)) // Default pprof routes served here
	}()

	log.Println(fmt.Sprintf("Starting server on port %s", config.Port))
	app.Listen(config.Port)
}
