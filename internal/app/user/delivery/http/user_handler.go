package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
)

type UserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserHandler(sserUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: sserUsecase,
	}
}

func (uh *UserHandler) Configure(r *mux.Router) {
	r.HandleFunc("/user/{nickname}/create", uh.UserCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{nickname}/profile", uh.UserGetProfileHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{nickname}/profile", uh.UserChangeProfileHandler).Methods(http.MethodPost)
}