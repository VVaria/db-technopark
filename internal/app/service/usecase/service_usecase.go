package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/service"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewServiceUsecase(serviceRepo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepo,
	}
}
