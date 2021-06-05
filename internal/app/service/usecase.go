package service

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ServiceUsecase interface {
	ClearDatabases() *errors.Error
	GetServiceStatus() (*models.Status, *errors.Error)
}
