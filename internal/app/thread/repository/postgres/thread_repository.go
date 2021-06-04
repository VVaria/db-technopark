package postgres

import (
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/models"
	"github.com/VVaria/db-technopark/internal/app/thread"
)

type ThreadRepository struct {
	conn *pgx.ConnPool
}

func NewThreadRepository(conn *pgx.ConnPool) thread.ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}