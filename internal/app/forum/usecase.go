package forum

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, *errors.Error)
	GetForumInfo(slug string) (*models.Forum, *errors.Error)
	GetForumUsers(slug string, params *models.Parameters) ([]*models.User, *errors.Error)
	GetForumThreads(slug string, params *models.Parameters) ([]*models.Thread, *errors.Error)
}
