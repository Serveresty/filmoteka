package main

import (
	"database/sql"
	"filmoteka/internal/database"
	"filmoteka/internal/transport"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../configs/.env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := dbInit()
	if err != nil {
		return fmt.Errorf("error with db connection:%v", err)
	}
	defer db.Close()
	database.DB = db

	err = database.CreateBaseTables()
	if err != nil {
		return fmt.Errorf("error while creating base tables:%v", err)
	}

	mux := http.NewServeMux()
	transport.Routes(mux)

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	addr := host + port

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
	return nil
}

func dbInit() (*sql.DB, error) {
	sqlDriver := os.Getenv("DATABASE_DRIVER")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbName)
	db, err := sql.Open(sqlDriver, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
