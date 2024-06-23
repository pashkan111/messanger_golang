package entities

// import (
// 	"context"
// 	"fmt"

// 	"messanger/src/entities"
// 	"messanger/src/repository/postgres_repos"

// 	"errors"

//	"github.com/jackc/pgconn"
//	"github.com/jackc/pgx/v4/pgxpool"
//	"github.com/sirupsen/logrus"
//
// )
type Chat struct {
	Id           int
	Name         string
	Participants []int
}
