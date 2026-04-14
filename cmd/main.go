package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failer to load .env")
	}
	connStr := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Can't to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Succses to connect to PoostgreSQL!")
}
