package postgres

import (
	"database/sql"
	"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/app/models"
)

type ForumRepository struct {
	conn *sql.DB
}

func NewForumRepository(conn *sql.DB) forum.ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}
