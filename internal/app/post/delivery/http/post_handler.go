package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/post"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
)

type PostHandler struct {
	postUsecase post.PostUsecase
}

func NewPostHandler(postUsecase post.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

func (ph *PostHandler) Configure(r *mux.Router) {
	r.HandleFunc("/post/{id}/details", ph.PostGetDetailsHanler).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/details", ph.PostChangeHandler).Methods(http.MethodPost)
}