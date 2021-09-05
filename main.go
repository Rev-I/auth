package main

import (
	"auth/controllers"
	dbutil "auth/db-util"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func init() {
	var err error
	// initialize connection pool to the database
	// this is a good way to handle concurrent database connection
	dbutil.Pool, err = pgxpool.Connect(context.Background(), dbutil.DB_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// initialize DB tables
	dbutil.InitTables()
}

func main() {
	// create a server using gin-gonic
	router := gin.Default()
	// register the routes
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
	router.GET("/me", controllers.Me)
	// start the server @ 127.0.0.1:8080
	router.Run()
}
