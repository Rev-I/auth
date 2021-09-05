package controllers

import (
	dbutil "auth/db-util"
	"auth/models"
	"auth/utilities"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user models.UserResource
	// set memory aside to read the request into
	bs, err := ioutil.ReadAll(c.Request.Body)
	// error handling if reading into a buffer fails
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	// parse json coming from the client into the user struct
	json.Unmarshal(bs, &user)
	// hash the user's password
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(password)
	// estblish a concurrent safe connection to the database
	conn, err := dbutil.Pool.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	defer conn.Release()
	var id string
	// register the user into the database
	err = conn.QueryRow(context.Background(),
		`
		INSERT INTO users (first_name, last_name, email, password, phone)
		values($1,$2,$3,$4,$5)
		RETURNING user_id
		`, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	user.ID = id
	token := utilities.GenerateToken(c, &user)
	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})

}

func Login(c *gin.Context) {
	var user models.UserResource
	var credentials models.UserCredentials
	// parse login details into 'credentials'
	c.ShouldBindJSON(&credentials)

	conn, err := dbutil.Pool.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	defer conn.Release()
	// find user with similar credentials
	err = conn.QueryRow(context.Background(), `
	SELECT * FROM users
	WHERE
	email =$1
	`, credentials.Email).Scan(
		&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// compare password from the request with the one saved in the DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failed"})
		c.Abort()
		return
	} else {
		token := utilities.GenerateToken(c, &user)
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token,
		})
	}

}

func Logout(c *gin.Context) {

	c.SetCookie(
		"jwt", "", int(time.Now().Add(-1*time.Hour).UnixNano()), "/", "localhost", false, true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})

}
func Me(c *gin.Context) {
	// retrieve the jwt token from the cookie
	jwtToken, err := c.Cookie("jwt")
	if jwtToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(utilities.SecretKey), nil
	})

	claims := token.Claims.(*jwt.StandardClaims)

	conn, err := dbutil.Pool.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	defer conn.Release()
	var user models.UserResource
	err = conn.QueryRow(context.Background(), `
	SELECT * FROM users
	WHERE
	user_id =$1
	`, &claims.Issuer).Scan(
		&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		c.Abort()
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
