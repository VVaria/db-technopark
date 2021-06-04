package postgres

import (
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/app/models"
)

type ForumRepository struct {
	conn *pgx.ConnPool
}

func NewForumRepository(conn *pgx.ConnPool) forum.ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}