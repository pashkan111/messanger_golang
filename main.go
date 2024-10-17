package main

import (
	"context"
	"fmt"
	"time"

	"messanger/src/dependencies"
	"messanger/src/services/event_broker"
	"messanger/src/views/api"
	"messanger/src/views/ws"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log := dependencies.GetLogger()
	postgres_pool := dependencies.GetPostgresPool(ctx, log)
	redis_pool := dependencies.GetRedisPool(ctx, log)
	broker := event_broker.RedisBroker{Client: redis_pool}

	router := mux.NewRouter()

	authHandler := api.NewAuthHandler(postgres_pool, log)
	router.HandleFunc("/auth/register", authHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.LoginUser).Methods("POST")

	wsHandler := ws.NewWSHandler(postgres_pool, log, &broker)
	router.HandleFunc("/process-events", wsHandler.HandleConnections)

	fmt.Println("Server is running on port 8080")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
