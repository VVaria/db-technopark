package postgres

import (
	"github.com/jackc/pgx"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/service"
)

type ServiceRepository struct {
	conn *pgx.ConnPool
}

func NewServiceRepository(conn *pgx.ConnPool) service.ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}
