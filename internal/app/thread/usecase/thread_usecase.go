package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/models"
	"github.com/VVaria/db-technopark/internal/app/thread"
	//"github.com/VVaria/db-technopark/internal/app/tools/uuid"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
	forumRepo  forum.ForumRepository
}

func NewThreadUsecase(threadRepo thread.ThreadRepository, forumRepo forum.ForumRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}
