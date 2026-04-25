package main

import (
	"auth/internals/db"
	authjwt "auth/internals/jwt"
	"auth/internals/store"
	"fmt"
	"os"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me-in-production"
	}

	db, err := db.New(dbURL, 10, 5, "5m")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	app := &application{
		config: config{
			address: ":8080",
		},
		store: store.NewStore(db),
		jwt:   authjwt.NewJWT(jwtSecret),
		dbConfig: dbConfig{
			address:      dbURL,
			maxOpenConns: 10,
			maxIdleConns: 5,
			maxIdleTime:  "5m",
		},
	}
	app.startServer()
}
