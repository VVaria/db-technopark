package forum

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, *errors.Error)
	GetForumInfo(slug string) (*models.Forum, *errors.Error)
	GetForumUsers(slug string, params *models.ForumUsersParameters) ([]*models.User, *errors.Error)
}
