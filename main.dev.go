//go:build dev

package main

import (
	"github.com/tnqbao/gau_blog_service/api/vote"
	"github.com/tnqbao/gau_blog_service/lib/ai"
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

	ai.LoadAPIKeys()
	router := routes.SetupRouter(db)
	go vote.StartSyncJob(db)
	router.Run(":8085")
}
