package utils

import (
	"fmt"
	"os"

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func GetPostgresPool(ctx context.Context, log *logrus.Logger) *pgxpool.Pool {
	godotenv.Load()

	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DATABASE"),
	)
	pool, err := pgxpool.Connect(ctx, postgresUrl)
	if err != nil {
		log.Error("Error with connecting to DB", err)
	}
	return pool
}
