package usecase

import (
	"strings"

	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/post"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/VVaria/db-technopark/internal/models"
)

type PostUsecase struct {
	postRepo post.PostRepository
	userRepo user.UserRepository
	forumRepo forum.ForumRepository
	threadRepo thread.ThreadRepository
}

func NewPostUsecase(postRepo post.PostRepository, userRepo user.UserRepository, forumRepo forum.ForumRepository, threadRepo thread.ThreadRepository) post.PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
		forumRepo: forumRepo,
		threadRepo: threadRepo,
	}
}

func (pu *PostUsecase) GetPostInfo(id int, related string) (interface{}, *errors.Error) {
	var allPostInf models.AllPostInfo
	var err error

	allPostInf.Post, err = pu.postRepo.SelectPostById(id)
	if err != nil {
		return nil, errors.Cause(errors.PostNotExist)
	}

	var items []string
	if strings.Contains(related, "user") {
		items = append(items, "user")
	}
	if strings.Contains(related, "forum") {
		items = append(items, "forum")
	}
	if strings.Contains(related, "thread") {
		items = append(items, "thread")
	}

	for _, i := range items {
		switch i {
		case "user":
			allPostInf.Author, err = pu.userRepo.SelectUserByNickname(allPostInf.Post.Author)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}
			break
		case "forum":
			allPostInf.Forum, err = pu.forumRepo.SelectForumBySlug(allPostInf.Post.Forum)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}
			break
		case "thread":
			allPostInf.Thread, err = pu.threadRepo.SelectThreadByID(allPostInf.Post.Thread)
			if err != nil {
				return nil, errors.UnexpectedInternal(err)
			}
			break
		}
	}

	var result interface{}
	result = allPostInf

	if allPostInf.Thread != nil {
		if models.IsUuid(allPostInf.Thread.Slug) {
			result = &models.AllPostInfoWithoutSlug{
				Post:   allPostInf.Post,
				Author: allPostInf.Author,
				Forum:  allPostInf.Forum,
				Thread: models.ThreadNoSlug(allPostInf.Thread),
			}
		}
	}
	return result, nil
}

func (pu *PostUsecase) ChangePostMessage(post *models.Post) (*models.Post, *errors.Error) {
	err := pu.postRepo.Update(post)
	if err != nil {
		return nil, errors.Cause(errors.PostNotExist)
	}

	return post, nil
}