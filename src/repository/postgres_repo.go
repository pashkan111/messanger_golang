package repository

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

func TestPostgres(ctx context.Context, conn *pgx.Conn, log *logrus.Logger) {
	log.Info("Testing Postgres connection")
	err := conn.Ping(ctx)
	if err != nil {
		log.Error("Error with pinging Postgres", err)
	}
	log.Info("Postgres connection is OK")
}
