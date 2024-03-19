package main

import (
	"richmedia/fravega/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/gofiber/fiber/v2"
)

func main() {
	app := gin.Default()

	app.Use(cors.Default())

	routes.Setup(app)

	app.Run(":8000")
}
