package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/VVaria/db-technopark/internal/models"
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

func (uh *UserHandler) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := &models.User{Nickname: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()

	userData, errE := uh.userUsecase.CreateUser(user)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONSuccess("Пользователь уже присутсвует в базе данных.", userData))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(errors.JSONSuccess(" Пользователь успешно создан.", user))
}

func (uh *UserHandler) UserGetProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, errE := uh.userUsecase.GetUserByNickname(nickname)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Информация о пользователе.", user))
}

func (uh *UserHandler) UserChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := &models.User{Nickname: nickname}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	defer r.Body.Close()

	errE := uh.userUsecase.UpdateProfile(user)
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Актуальная информация о пользователе после изменения профиля.", user))
}
