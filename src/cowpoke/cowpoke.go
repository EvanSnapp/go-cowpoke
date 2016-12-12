package main

import (
	"middleware"

	"os"

	"routes"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			panic("a .env file is required for local development")
		}
	}

	r := gin.Default()
	//global middleware
	r.Use(middleware.Authenticate())

	//wireup all routes
	api := r.Group("/api")
	{
		routes.AddStatusRoutes(api)
		routes.AddStackRoutes(api)
	}

	//start the server
	r.Run(fmt.Sprintf(":%s", os.Getenv("HOST_PORT")))
}
