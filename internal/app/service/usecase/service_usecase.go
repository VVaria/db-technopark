package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/service"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewServiceUsecase(serviceRepo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepo,
	}
}

func (su *ServiceUsecase) ClearDatabases() *errors.Error {
	err := su.serviceRepo.ClearAll()
	if err != nil {
		return errors.UnexpectedInternal(err)
	}
	return nil
}

func (su *ServiceUsecase) GetServiceStatus() (*models.Status, *errors.Error) {
	status, err := su.serviceRepo.ServiceStatus()
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	return status, nil
}
