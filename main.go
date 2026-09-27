package main

import (
	"h-lab/database"
	"h-lab/framework/fiber"
	"h-lab/router"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	godotenv.Load()

	// เชื่อมต่อ Database และทำ Migration
	database.Connect()

	r := fiber.NewFiberApp()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API
	router.SetRouter(r)
	if err := r.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
