package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
	forumRepo  forum.ForumRepository
	userRepo user.UserRepository
}

func NewThreadUsecase(threadRepo thread.ThreadRepository, forumRepo forum.ForumRepository, userRepo user.UserRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo: userRepo,
	}
}
func (tu *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	_, err := tu.forumRepo.SelectForumBySlug(thread.Forum)
	if err != nil {
		return nil, errors.Cause(errors.ForumNotExist)
	}

	_, err = tu.userRepo.SelectUserByNickname(thread.Author)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	if thread.Slug == "" {
		slug, err := uuid.NewRandom()
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
		thread.Slug = slug.String()
	}


	err = tu.threadRepo.InsertThread(thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.ThreadNotExist)
			}

			if pgErr.Code == "23505" {
				threadInfo, err := tu.threadRepo.SelectThreadByID(thread.Id)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}
				return threadInfo, errors.Cause(errors.ForumCreateThreadConflict)
			}
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return thread, nil
}