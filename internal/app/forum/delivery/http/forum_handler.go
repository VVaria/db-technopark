package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ForumHandler struct {
	forumUsecase forum.ForumUsecase
	threadUsecase thread.ThreadUsecase
}

func NewForumHandler(forumUsecase forum.ForumUsecase, threadUsecase thread.ThreadUsecase) *ForumHandler {
	return &ForumHandler{
		forumUsecase: forumUsecase,
		threadUsecase: threadUsecase,
	}
}

func (fh *ForumHandler) Configure(r *mux.Router) {
	r.HandleFunc("/forum/create", fh.ForumCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/details", fh.ForumDetailsHandler).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/create", fh.ForumCreateThreadHandler).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/users", fh.ForumUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/threads", fh.ForumThreadsHandler).Methods(http.MethodGet)
}

func (fh *ForumHandler) ForumCreateHandler(w http.ResponseWriter, r *http.Request) {
	forum := &models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		return
	}

	forumInfo, errE := fh.forumUsecase.CreateForum(forum)
	if errE.ErrorCode == errors.UserNotExist {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONMessage("Владелец форума не найден."))
		return
	}
	if errE.ErrorCode == errors.ForumCreateConflict {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONMessage("Форум уже присутсвует в базе данных.", forumInfo))
		return
	}
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Форум успешно создан.", forumInfo))
}


func (fh *ForumHandler) ForumDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	forumInfo, errE := fh.forumUsecase.GetForumInfo(slug)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о форуме.", forumInfo))
}

func (fh *ForumHandler) ForumCreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	thread := &models.Thread{Forum: slug}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		return
	}
	flag := thread.Slug == ""

	threadInfo, errE := fh.threadUsecase.CreateThread(thread)
	if errE.ErrorCode == errors.ForumCreateThreadConflict {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONMessage("Ветка обсуждения уже присутсвует в базе данных.", threadInfo))
		return
	}
	if errE.ErrorCode == errors.UserNotExist || errE.ErrorCode == errors.ForumNotExist {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONMessage("Автор ветки или форум не найдены."))
		return
	}
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	if flag {
		threadWithout := models.ThreadNoSlug(threadInfo)
		w.WriteHeader(http.StatusOK)
		w.Write(errors.JSONMessage(" Ветка обсуждения успешно создана.", threadWithout))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage(" Ветка обсуждения успешно создана.", threadInfo))
}

func (fh *ForumHandler) ForumUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	params := &models.ForumUsersParameters{}
	var err error
	params.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		params.Limit = 100
	}
	params.Since = r.URL.Query().Get("since")
	params.Desc, err = strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		params.Desc = true
	}

	users, errE := fh.forumUsecase.GetForumUsers(slug, params)
	if errE.ErrorCode == errors.ForumNotExist {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONMessage("Форум отсутсвует в системе."))
		return
	}
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONMessage("Информация о пользователях форума.", users))
}