package postgres

import (
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/models"
	"github.com/VVaria/db-technopark/internal/app/user"
)

type UserRepository struct {
	conn *pgx.ConnPool
}

func NewUserRepository(conn *pgx.ConnPool) user.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}
