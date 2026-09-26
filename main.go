package main

import (
	"h-lab/database"
	app "h-lab/framework/fiber"

	"h-lab/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	godotenv.Load()

	// เชื่อมต่อ Database และทำ Migration
	database.Connect()

	r := app.NewFiberApp()
	//api
	router.SetRouter(r)
	if err := r.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
