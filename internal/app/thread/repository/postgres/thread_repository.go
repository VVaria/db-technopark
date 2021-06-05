package postgres

import (
	"database/sql"
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/thread"
)

type ThreadRepository struct {
	conn *sql.DB
}

func NewThreadRepository(conn *sql.DB) thread.ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}
