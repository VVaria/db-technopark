package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/post"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strconv"

	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
	forumRepo  forum.ForumRepository
	userRepo user.UserRepository
	postRepo post.PostRepository
}

func NewThreadUsecase(threadRepo thread.ThreadRepository, forumRepo forum.ForumRepository, userRepo user.UserRepository, postRepo post.PostRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo: userRepo,
		postRepo: postRepo,
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

func (tu *ThreadUsecase) CreateThreadPosts(thread string, posts []*models.Post) ([]*models.Post, *errors.Error) {
	threadInfo, errE := tu.GetThreadInfo(thread)
	if errE != nil {
		return nil, errors.Cause(errors.ThreadNotExist)
	}
	postsInfo, err := tu.postRepo.InsertPosts(posts, threadInfo.Id, threadInfo.Forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "12345" {
				return nil, errors.Cause(errors.PostWrongThread)
			}

			if pgErr.Code == "23503" {
				return nil, errors.Cause(errors.UserNotExist)
			}
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return postsInfo, nil
}

func (tu *ThreadUsecase) GetThreadInfo(thread string) (*models.Thread, *errors.Error) {
	var threadInfo *models.Thread
	id, err := strconv.Atoi(thread)
	if err != nil {
		threadInfo, err = tu.threadRepo.SelectThreadBySlug(thread)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	} else {
		threadInfo, err = tu.threadRepo.SelectThreadByID(id)
		if err != nil {
			return nil, errors.Cause(errors.ThreadNotExist)
		}
	}
	return threadInfo, nil
}

func (tu *ThreadUsecase) RefreshThread(threadId string, thread *models.Thread) (*models.Thread, *errors.Error) {
	id, err := strconv.Atoi(threadId)
	if err != nil {
		thread.Slug = threadId
		thread.Id = 0
	} else {
		thread.Slug = ""
		thread.Id = id
	}

	err = tu.threadRepo.UpdateThread(thread)
	if err != nil {
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	return thread, nil
}