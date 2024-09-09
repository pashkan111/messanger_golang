package event_handlers

import (
	"context"
	"messanger/src/events/request_events"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type EventHandler func(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.RequestEventInterface,
) (interface{}, error)
