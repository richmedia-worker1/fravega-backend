package routes

import (
	"richmedia/fravega/internal/api/auth"
	"richmedia/fravega/internal/api/cataloge"
	"richmedia/fravega/internal/api/user"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Setup(app *gin.Engine) {
	api := app.Group("v1")
	api.Use(CORSMiddleware())
	auth.Setup(api)
	user.Setup(api)
	cataloge.Setup(api)
}
