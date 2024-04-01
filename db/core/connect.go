package core

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {
	connectionString := createConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("Error opening database connection")
	}
	return db
}
func createConnectionString() string {
	dbName, dbUser, dbPassword, dbHost, dbPort := readEnvVariables()
	result := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	slog.Info(fmt.Sprintf("Created postgresql path: %v", result))
	return result
}
func readEnvVariables() (dbName, dbUser, dbPassword, dbHost, dbPort string) {
	if err := godotenv.Load(); err != nil {
		slog.Error("error reading environment variables")
	}
	dbName, exists := os.LookupEnv("DATABASE_NAME")
	handleExists(exists)
	dbUser, exists = os.LookupEnv("DATABASE_USER")
	handleExists(exists)
	dbPassword, exists = os.LookupEnv("DATABASE_PASSWORD")
	handleExists(exists)
	dbHost, exists = os.LookupEnv("DATABASE_HOST")
	handleExists(exists)
	dbPort, exists = os.LookupEnv("DATABASE_PORT")
	handleExists(exists)
	return dbName, dbUser, dbPassword, dbHost, dbPort
}
func handleExists(exists bool) {
	if !exists {
		slog.Error("No such env variable")
	}
}
