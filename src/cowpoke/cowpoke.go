package main

import (
	"log"
	"middleware/authentication"
	"middleware/errors"
	"os"
	"routes"

	"configuration"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	if errs := configuration.Init(); errs != nil {
		msg := fmt.Sprintf("service could not start due to the following configuration errors:\n %s", errs)
		panic(msg)
	}

	r := gin.Default()
	//global middleware
	r.Use(authentication.Authenticate())
	r.Use(errors.HandlePublicError())

	//wireup all routes
	api := r.Group("/api")
	{
		routes.AddStatusRoutes(api)
		routes.AddStackRoutes(api)
	}

	listenAddr := fmt.Sprintf(":%s", os.Getenv("HOST_PORT"))
	log.Printf("service running on %s\n", listenAddr)
	//start the server
	r.Run(listenAddr)
}
