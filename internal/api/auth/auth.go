package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	Id       int
	Login    string `json:"login"`
	Password string `json:"password"`
}

var Users []Auth
var secretKey = []byte("my-secret-key")

func Setup(rg *gin.RouterGroup) {

	Users = []Auth{
		{Id: 0, Login: "test", Password: "123"},
		{Id: 1, Login: "test@g.co", Password: "123"},
	}

	api := rg.Group("auth")
	api.GET("", WithAuth, getGet)
	api.GET("/signin", getSignin)
	api.POST("/signup", postSignup)
}

func getGet(c *gin.Context) {
	c.String(http.StatusOK, "test, claims:\n", c.MustGet("claims").(jwt.MapClaims))
}
func getSignin(c *gin.Context) {
	passwd := c.Request.URL.Query().Get("password")
	login := c.Request.URL.Query().Get("login")
	if passwd == "" || login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login or password can't be null"})
		return
	}
	log.Println("login: ", login)
	log.Println("passwd: ", passwd)

	for _, user := range Users {
		if user.Login == login && user.Password == passwd {
			token := createToken(user)
			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this login and password not found!"})
}

func postSignup(c *gin.Context) {
	var creds Auth
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil || (creds == Auth{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credencials can't be null"})
		return
	}
	if creds.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login can't be null"})
		return
	}
	if creds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password can't be null"})
		return
	}
	for _, user := range Users {
		if user.Login == creds.Login {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user with this login is already created"})
			return
		}
	}
	user := Auth{
		Id:       Users[len(Users)-1].Id + 1,
		Login:    creds.Login,
		Password: creds.Password,
	}
	Users = append(Users, user)

	token := createToken(user)
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Создание JWT-токена
func createToken(user Auth) string {
	claims := jwt.MapClaims{
		"login": user.Login,
		"id":    user.Id,
		"exp":   time.Now().Add(time.Hour * 24 * 180).Unix(), // Токен действителен 180 дней(пол года)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)

	return signedToken
}

// Middleware для проверки токена
func WithAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorizated"})
		c.Abort()
		return
	}

	// Проверяем, что токен начинается с "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
		c.Abort()
		return
	}

	// Извлекаем сам токен
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		c.Abort()
		return
	}

	// достаем значения
	claims := token.Claims.(jwt.MapClaims)
	log.Println(claims)
	for _, user := range Users {
		if user.Id == int(claims["id"].(float64)) {
			c.Set("claims", claims)
			c.Next()
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this token not found"})
	c.Abort()
}
