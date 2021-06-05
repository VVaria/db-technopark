package user

import (
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type UserUsecase interface {
	CreateUser(user *models.User) ([]models.User, *errors.Error)
	GetUserByNickname(nickname string) (*models.User, *errors.Error)
	UpdateProfile(user *models.User) *errors.Error
}
