package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/post"
	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
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