package routes

import (
	"richmedia/fravega/internal/api/auth"
	"richmedia/fravega/internal/api/cataloge"
	"richmedia/fravega/internal/api/user"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	api := app.Group("v1")
	auth.Setup(api)
	user.Setup(api)
	cataloge.Setup(api)
}
