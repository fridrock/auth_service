package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CreateConnection() *sqlx.DB {
	connectionString := createConnectionString()
	db, err := sqlx.Connect("postgres", connectionString)
	//db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("Error opening database connection")
	}
	slog.Info("Created postgresql connection")
	return db
}
func createConnectionString() string {
	dbName, dbUser, dbPassword, dbHost, dbPort := readEnvVariables()
	result := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
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
