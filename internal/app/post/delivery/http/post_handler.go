package http

import (
	"encoding/json"
	"fmt"
	"github.com/VVaria/db-technopark/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	//"github.com/VVaria/db-technopark/internal/app/models"
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


func (ph *PostHandler) PostGetDetailsHanler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]

	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}

	related := r.URL.Query().Get("related")
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


	postInfo, errE := ph.postUsecase.GetPostInfo(id, related)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Информация о ветке обсуждения.", postInfo))
}