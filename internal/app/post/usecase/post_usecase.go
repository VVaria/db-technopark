package usecase

import (
	"strconv"

	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/models"
	"github.com/VVaria/db-technopark/internal/app/post"
	//"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
)

type PostUsecase struct {
	postRepo   post.PostRepository
	//threadRepo thread.ThreadRepository
}

func NewPostUsecase(postRepo post.PostRepository) post.PostUsecase {
	return &PostUsecase{
		postRepo:   postRepo,
	}
}
