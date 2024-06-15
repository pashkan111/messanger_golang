package utils

import (
	"os"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func GetPostgresConn(log *logrus.Logger) *pgx.Conn {
	godotenv.Load()

	port, _ := strconv.Atoi(os.Getenv("PG_POST"))
	postgresUrl := pgx.ConnConfig{
		Host:     os.Getenv("PG_HOST"),
		Port:     uint16(port),
		Database: os.Getenv("PG_DATABASE"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
	}
	conn, err := pgx.Connect(postgresUrl)
	if err != nil {
		log.Error("Error with connecting to DB", err)
	}
	return conn
}
