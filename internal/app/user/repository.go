package user

import "github.com/VVaria/db-technopark/internal/models"

type UserRepository interface {
	SelectUsers(user *models.User) ([]models.User, error)
	InsertUser(user *models.User) error
	SelectUserByNickname(nickname string) (*models.User, error)
	Update(user *models.User) error
	SelectForumUsers(slug string, params *models.Parameters) ([]*models.User, error)
}
