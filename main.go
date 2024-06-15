package main

import (
	"context"
	"time"

	"messanger/src/repository"
	"messanger/src/utils"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log := utils.GetLogger()
	postgres_conn := utils.GetPostgresConn(log)
	repository.TestPostgres(ctx, postgres_conn, log)
}
