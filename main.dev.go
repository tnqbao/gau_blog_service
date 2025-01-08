//go:build dev

package main

import (
	"github.com/tnqbao/gau_blog_service/api/vote"
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.InitRedis()
	db := config.InitDB()
	router := routes.SetupRouter(db)
	go vote.StartSyncJob(db)
	router.Run(":8085")
}
