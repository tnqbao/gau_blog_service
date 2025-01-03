package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/routes"
)

func main() {
	err := godotenv.Load("/gau_blog/.env.blog")
	if err != nil {
		log.Fatalf("Error loading .env file: ", err)
	}
	db := config.InitDB()
	router := routes.SetupRouter(db)
	router.Run(":8085")
}
