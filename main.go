package main

import (
	"log"
	"os"

	"github.com/itzblinkzy/act-backend/config"
	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/external"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.InitDB()
	defer db.Close()
	e := config.InitEcho()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	external.InitResty()
	config.InitRoutes(e)
	e.Logger.Fatal(e.Start(":" + port))
}
