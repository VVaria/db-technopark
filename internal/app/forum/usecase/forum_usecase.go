package usecase

import (
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
)

type ForumUsecase struct {
	forumRepo forum.ForumRepository
}

func NewForumUsecase(forumRepo forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRepo: forumRepo,
	}
}