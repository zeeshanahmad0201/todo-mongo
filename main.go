package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/zeeshanahmad0201/todo-mongo/common"
	"github.com/zeeshanahmad0201/todo-mongo/database"
	"github.com/zeeshanahmad0201/todo-mongo/router"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.InitMongoDB()
	defer database.CloseMongo()

	port := "8080"
	addr := ":" + port
	r := router.InitRouter()

	log.Printf("Starting server at %s...\n", port)

	if err := http.ListenAndServe(addr, r); err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Exit: true,
		})
	}
}
