package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"
)

type ForumUsecase struct {
	forumRepo forum.ForumRepository
	userRepo user.UserRepository
	threadRepo thread.ThreadRepository
}

func NewForumUsecase(forumRepo forum.ForumRepository, userRepo user.UserRepository, threadRepo thread.ThreadRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRepo: forumRepo,
		userRepo: userRepo,
		threadRepo: threadRepo,
	}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, *errors.Error) {
	userInfo, err := fu.userRepo.SelectUserByNickname(forum.User)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}
	forum.User = userInfo.Nickname
	err = fu.forumRepo.CreateForum(forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return nil, errors.Cause(errors.ForumNotExist)
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			forumModel, err := fu.forumRepo.SelectForumBySlug(forum.Slug)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}
			return forumModel, errors.Cause(errors.ForumCreateConflict)
		}
		return nil, errors.UnexpectedInternal(err)
	}
	return forum, nil
}


func (fu *ForumUsecase) GetForumInfo(slug string) (*models.Forum, *errors.Error) {
	forum, err := fu.forumRepo.SelectForumBySlug(slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Cause(errors.ForumNotExist)
		}
		return nil, errors.UnexpectedInternal(err)
	}
	return forum, nil
}


func (fu *ForumUsecase) GetForumUsers(slug string, params *models.Parameters) ([]*models.User, *errors.Error) {
	_, err := fu.forumRepo.SelectForumBySlug(slug)
	if err != nil {
		return nil, errors.Cause(errors.ForumNotExist)
	}

	users, err := fu.userRepo.SelectForumUsers(slug, params)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return users, nil
}

func (fu *ForumUsecase) GetForumThreads(slug string, params *models.Parameters) ([]*models.Thread, *errors.Error) {
	_, err := fu.forumRepo.SelectForumBySlug(slug)
	if err != nil {
		return nil, errors.Cause(errors.ForumNotExist)
	}

	threads, err := fu.threadRepo.SelectForumThreads(slug, params)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return threads, nil
}