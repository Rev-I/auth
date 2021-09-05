package utilities

import (
	"auth/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// secret key for our app, it should be hidden in environment variables
const SecretKey = "secret"

// this function generate JWT when a user registers or when they log in
func GenerateToken(c *gin.Context, user *models.UserResource) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Issuer:    user.ID,
		ExpiresAt: jwt.NewTime(float64(time.Now().Add(24 * time.Hour).UnixNano())),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to authonticate"})
		return ""
	}
	c.SetCookie(
		"jwt", token, int(time.Now().Add(24*time.Hour).UnixNano()), "/", "localhost", false, true,
	)
	return token
}
