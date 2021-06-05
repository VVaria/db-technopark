package http

import (
	"encoding/json"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

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
	r.HandleFunc("/post/{id}/details", ph.PostGetDetailsHandler).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/details", ph.PostChangeHandler).Methods(http.MethodPost)
}


func (ph *PostHandler) PostGetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]

	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}

	related := r.URL.Query().Get("related")

	postInfo, errE := ph.postUsecase.GetPostInfo(id, related)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о ветке обсуждения.", postInfo))
}

func (ph *PostHandler) PostChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]

	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}

	postInfo := &models.Post{ID: id}
	err = json.NewDecoder(r.Body).Decode(&postInfo)
	if err != nil {
		return
	}

	postInfo, errE := ph.postUsecase.ChangePostMessage(postInfo)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация об измененном сообщении.", postInfo))
}
