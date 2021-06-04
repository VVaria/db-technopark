package usecase

import (
	"strconv"

	"github.com/jackc/pgx"

	//"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
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