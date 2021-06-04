package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/VVaria/db-technopark/internal/app/forum"
	//"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	//"github.com/VVaria/db-technopark/internal/app/models"
)

type ForumHandler struct {
	forumUsecase forum.ForumUsecase
}

func NewForumHandler(forumUsecase forum.ForumUsecase) *ForumHandler {
	return &ForumHandler{
		forumUsecase: forumUsecase,
	}
}

func (fh *ForumHandler) Configure(r *mux.Router) {
	r.HandleFunc("/forum/create", fh.ForumCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/details", fh.ForumDetailsHandler).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/create", fh.ForumCreateThreadHandler).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/users", fh.ForumUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/forum/{slug}/threads", fh.ForumThreadsHandler).Methods(http.MethodGet)
}
