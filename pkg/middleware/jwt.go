package middleware

import (
	"net/http"
	"strings"

	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

// Claims for JWT. We store all the JWT default claims + username for this
// software.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Check the authorization from the header bearer token. If the authorization is
// good does nothing, else it aborts the Gin context.
func Auth(c *gin.Context) {
	config, err := config.GetConfig()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()

		return
	}

	key := []byte(config.String("jwt.token"))
	bearer := c.Request.Header.Get("Authorization")
	claims := &Claims{}

	// If header does not start with "Bearer " better to stop here
	if !strings.HasPrefix(bearer, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})

		c.Abort()
		return
	}

	// JWT is parsed only by the last part of the Authorization header
	token, err := jwt.ParseWithClaims(strings.Split(bearer, " ")[1], claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		}

		c.Abort()
		return
	} else if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		c.Abort()

		return
	}
}
