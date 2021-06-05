package postgres

import (
	"database/sql"
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/post"
	//"github.com/VVaria/db-technopark/internal/app/tools/null"
)

type PostRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.PostRepository {
	return &PostRepository{
		conn: conn,
	}
}
