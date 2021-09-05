package dbutil

import (
	"auth/models"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const DB_URL = "postgres://postgres:Khumalo87@localhost:5432/auth"

var Pool *pgxpool.Pool

func InitTables() {
	// connect to the database
	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	// install a plugin to use uuid
	conn.Exec(context.Background(), `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`)
	stmt, err := conn.Prepare(context.Background(), "createUser", models.User)
	if err != nil {
		fmt.Println("I ran?")
		fmt.Println(err)
		os.Exit(1)
	}

	conn.Exec(context.Background(), stmt.Name)
}
