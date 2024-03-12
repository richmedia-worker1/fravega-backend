package user

import (
	"net/http"
	"richmedia/fravega/internal/api/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Setup(rg *gin.RouterGroup) {
	api := rg.Group("user")
	api.GET("", auth.WithAuth, getGetUser)

}

func getGetUser(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	var user auth.Auth
	for _, _user := range auth.Users {
		if _user.Id == int(claims["id"].(float64)) {
			user = _user
		}
	}
	if (user == auth.Auth{}) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "claims is empty or user can't be extracted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
