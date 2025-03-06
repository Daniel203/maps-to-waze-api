package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentVariables()

    initLogger()
	initRouter()
}

func loadEnvironmentVariables() {
	env, isPresent := os.LookupEnv("ENV")

	// If production do nothing
	if isPresent && env == "prod" {
		return
	}

	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("Failed to load the environment variables: %w", err))
	}
}

