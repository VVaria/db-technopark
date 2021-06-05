package post

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type PostUsecase interface {
	GetPostInfo(id int, related string) (interface{}, *errors.Error)
	ChangePostMessage(post *models.Post) (*models.Post, *errors.Error)
}
