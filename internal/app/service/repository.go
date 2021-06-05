package service

import "github.com/VVaria/db-technopark/internal/models"

type ServiceRepository interface {
	ClearAll() error
	ServiceStatus() (*models.Status, error)
}
