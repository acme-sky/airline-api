package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/acme-sky/airline-api/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Struct used for login
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Generate a valid JWT with HMAC 256 with an username and an expiration time of
// 1 hour. Key is stored in env.
func generateJWT() (string, error) {
	key := []byte(config.GetConfig().String("jwt.token"))
	expiration := time.Now().Add(time.Hour)
	claims := &middleware.Claims{
		Username: "username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)

}

// Handler used to login the system and get a JWT to make requests.
// Password is stored as SHA256 hashed in database.
func LoginHandler(c *gin.Context) {
	db := db.GetDb()

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(input.Password)))

	if err := db.Where("username = ? and password = ?", input.Username, password).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	token, _ := generateJWT()

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
