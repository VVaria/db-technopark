package postgres

import (
	"github.com/VVaria/db-technopark/internal/models"
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

func (sr *ServiceRepository) ClearAll() error {
	_, err := sr.conn.Exec(`TRUNCATE users, threads, forum, posts, votes, forum_users`)
	return err
}

func (sr *ServiceRepository) ServiceStatus() (*models.Status, error) {
	var status models.Status
	err := sr.conn.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&status.User)
	if err != nil {
		status.User = 0
		return nil, err
	}
	err = sr.conn.QueryRow(`SELECT COUNT(*) FROM forum;`).Scan(&status.Forum)
	if err != nil {
		status.Forum = 0
		return nil, err
	}
	err = sr.conn.QueryRow(`SELECT COUNT(*) FROM thread;`).Scan(&status.Thread)
	if err != nil {
		status.Thread = 0
		return nil, err
	}
	err = sr.conn.QueryRow(`SELECT COUNT(*) FROM post;`).Scan(&status.Post)
	if err != nil {
		status.Post = 0
		return nil, err
	}
	return &status, nil
}