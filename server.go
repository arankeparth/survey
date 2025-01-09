package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"survey-service/config"
	"survey-service/routes"
)

func main() {
	ctx := context.Background()
	routes.SetRoutes(ctx)
	
	//config.InitializeProcess()

	go func() {
		log.Println(fmt.Sprintf("Starting server on port %s", config.PprofPort))
		log.Println(http.ListenAndServe(":6060", nil)) // Default pprof routes served here
	}()

	log.Println(fmt.Sprintf("Starting server on port %s", config.Port))
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
