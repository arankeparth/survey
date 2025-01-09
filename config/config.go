package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
)

var (
	// Port is the port number the server will listen to.
	Port string
	// PProf is the port number for pprof.
	PprofPort string
	// DBHost is the host of the database.
	DBHost string
)

// Init function to load .env file and grab values in it.
func initializeProcess() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(err.Error())
	}
	loadEnv()
}

func init() {
	initializeProcess()
}

func loadEnv() {
	Port = os.Getenv("PORT")
	DBHost = os.Getenv("DB_HOST")
	PprofPort = os.Getenv("PPROF_PORT")
}
