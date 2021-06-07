package thread

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, *errors.Error)
	CreateThreadPosts(thread *models.Thread, posts []*models.Post) ([]*models.Post, *errors.Error)
	GetThreadInfo(thread string) (*models.Thread, *errors.Error)
	RefreshThread(threadId string, thread *models.Thread) (*models.Thread, *errors.Error)
	GetThreadPosts(slug string, params *models.ThreadPostParameters) ([]*models.Post, *errors.Error)
	ThreadVote(slug string, vote *models.Vote) (*models.Thread, *errors.Error)
}
