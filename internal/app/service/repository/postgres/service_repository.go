package postgres

import (
	"database/sql"
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/service"
)

type ServiceRepository struct {
	conn *sql.DB
}

func NewServiceRepository(conn *sql.DB) service.ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}
