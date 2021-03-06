package usecase

import (
	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/post"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strconv"
)

type ThreadUsecase struct {
	threadRepo thread.ThreadRepository
	forumRepo  forum.ForumRepository
	userRepo   user.UserRepository
	postRepo   post.PostRepository
}

func NewThreadUsecase(threadRepo thread.ThreadRepository, forumRepo forum.ForumRepository, userRepo user.UserRepository, postRepo post.PostRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo:   userRepo,
		postRepo:   postRepo,
	}
}
func (tu *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	forum, err := tu.forumRepo.SelectForumBySlug(thread.Forum)
	if err != nil {
		return nil, errors.Cause(errors.ForumNotExist)
	}

	auth, err := tu.userRepo.SelectUserByNickname(thread.Author)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	thread.Forum = forum.Slug
	thread.Author = auth.Nickname
	if thread.Slug == "" {
		slug, err := uuid.NewRandom()
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
		thread.Slug = slug.String()
	}

	err = tu.threadRepo.InsertThread(thread)
	if err != nil {
		threadInfo, err := tu.threadRepo.SelectThreadBySlug(thread.Slug)
		if err != nil {
			return nil, errors.Cause(errors.ForumNotExist)
		}
		return threadInfo, errors.Cause(errors.ForumCreateThreadConflict)
	}

return thread, nil
}

func (tu *ThreadUsecase) CreateThreadPosts(thread *models.Thread, posts []*models.Post) ([]*models.Post, *errors.Error) {
	postsInfo, err := tu.postRepo.InsertPosts(posts, thread.Id, thread.Forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "00409" {
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

func (tu *ThreadUsecase) GetThreadPosts(slug string, params *models.ThreadPostParameters) ([]*models.Post, *errors.Error) {
	threadInfo, errE := tu.GetThreadInfo(slug)
	if errE != nil {
		return nil, errors.Cause(errors.ThreadNotExist)
	}

	var posts []*models.Post
	var err error
	switch params.Sort {
	case "flat":
		posts, err = tu.postRepo.SelectThreadPostsFlat(threadInfo.Id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
		break
	case "tree":
		posts, err = tu.postRepo.SelectThreadPostsTree(threadInfo.Id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
		break
	case "parent_tree":
		posts, err = tu.postRepo.SelectThreadPostsParent(threadInfo.Id, params)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
		break
	}

	return posts, nil
}

func (tu *ThreadUsecase) ThreadVote(slug string, vote *models.Vote) (*models.Thread, *errors.Error) {
	threadInfo, errE := tu.GetThreadInfo(slug)
	if errE != nil {
		return nil, errE
	}
	vote.Thread = threadInfo.Id

	err := tu.threadRepo.InsertVote(vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return nil, errors.Cause(errors.UserNotExist)
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			err = tu.threadRepo.UpdateVote(vote)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}

			threadInfo, errE = tu.GetThreadInfo(slug)
			if errE != nil {
				return nil, errE
			}
			return threadInfo, nil
		}
		return nil, errors.UnexpectedInternal(err)
	}
	threadInfo, errE = tu.GetThreadInfo(slug)
	if errE != nil {
		return nil, errE
	}
	return threadInfo, nil
}
