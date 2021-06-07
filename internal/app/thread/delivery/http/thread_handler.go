package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadHandler struct {
	threadUsecase thread.ThreadUsecase
}

func NewThreadHandler(threadUsecase thread.ThreadUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUsecase: threadUsecase,
	}
}

func (th *ThreadHandler) Configure(r *mux.Router) {
	r.HandleFunc("/thread/{slug_or_id}/create", th.ThreadCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/details", th.ThreadGetDetailsHandler).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/details", th.ThreadRefreshHandler).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/posts", th.ThreadPostsHandler).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/vote", th.ThreadVoteHandler).Methods(http.MethodPost)
}

func (th *ThreadHandler) ThreadCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thread := vars["slug_or_id"]

	var posts []*models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		return
	}

	t, errE := th.threadUsecase.GetThreadInfo(thread)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	if len(posts) == 0 {
		body := []models.Post{}
		w.WriteHeader(http.StatusCreated)
		w.Write(errors.JSONMessage("Информация о ветке обсуждения.", body))
		return
	}

	postsInfo, errE := th.threadUsecase.CreateThreadPosts(t, posts)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	var result []*models.PostShort
	for _, i := range postsInfo {
		p := &models.PostShort{
			ID: i.ID,
			Thread: i.Thread,
			Author: i.Author,
			Message: i.Message,
			Created: i.Created,
			Forum: i.Forum,
			Parent: i.Parent,
		}
		result = append(result, p)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(errors.JSONMessage("Посты успешно созданы.", result))
}

func (th *ThreadHandler) ThreadGetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thread := vars["slug_or_id"]

	threadInfo, errE := th.threadUsecase.GetThreadInfo(thread)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	if models.IsUuid(threadInfo.Slug) {
		var body interface{}
		body = models.ThreadNoSlug(threadInfo)
		w.WriteHeader(http.StatusOK)
		w.Write(errors.JSONMessage("Информация о ветке обсуждения.", body))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о ветке обсуждения.", threadInfo))
}

func (th *ThreadHandler) ThreadRefreshHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thread := vars["slug_or_id"]

	data := &models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return
	}

	threadInfo, errE := th.threadUsecase.RefreshThread(thread, data)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	var body interface{}
	body = threadInfo
	if models.IsUuid(threadInfo.Slug) {
		body = models.ThreadNoSlug(threadInfo)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о ветке обсуждения.", body))
}

func (th *ThreadHandler) ThreadPostsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thread := vars["slug_or_id"]

	params := &models.ThreadPostParameters{}
	var err error
	params.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		params.Limit = 100
	}
	params.Since, err = strconv.Atoi(r.URL.Query().Get("since"))
	if err != nil {
		params.Since = 0
	}
	params.Desc, err = strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		params.Desc = false
	}
	params.Sort = r.URL.Query().Get("sort")
	if params.Sort == "" {
		params.Sort = "flat"
	}

	posts, errE := th.threadUsecase.GetThreadPosts(thread, params)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}


	if len(posts) == 0 {
		body := []models.Post{}
		w.WriteHeader(http.StatusOK)
		w.Write(errors.JSONMessage("Информация о ветке обсуждения.", body))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о сообщениях форума.", posts))
}

func (th *ThreadHandler) ThreadVoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thread := vars["slug_or_id"]

	vote := &models.Vote{}
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		return
	}

	threadInfo, errE := th.threadUsecase.ThreadVote(thread, vote)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	var body interface{}
	body = threadInfo
	if models.IsUuid(threadInfo.Slug) {
		body = models.ThreadNoSlug(threadInfo)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о ветке обсуждения.", body))
}
