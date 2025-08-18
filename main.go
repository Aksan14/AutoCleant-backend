package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"reset/config"
	"reset/routes"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	appPort := os.Getenv("APP_PORT")
	fmt.Println(" http://localhost:" + appPort)

	db, err := config.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	routes.Routes(db, appPort)
}