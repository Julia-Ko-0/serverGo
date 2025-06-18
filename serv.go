package main

import (
	"log"
	"os"
	"serverGo/db"
	"serverGo/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		log.Println(" JWT_SECRET is EMPTY! Please check your .env file or environment variables.")
	}

	db.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Разрешаем запросы только с фронта на этом домене
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Если нужны cookies
		MaxAge:           12 * time.Hour,
	}))
	// r.Use(cors.Default())

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
