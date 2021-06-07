package usecase

import (
	"github.com/jackc/pgx"

	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/app/user"
	"github.com/VVaria/db-technopark/internal/models"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uu *UserUsecase) CreateUser(user *models.User) ([]models.User, *errors.Error)  {
	var users []models.User
	users, err := uu.userRepo.SelectUsers(user)
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.UnexpectedInternal(err)
	}
	if len(users) != 0 {
		return users, errors.Cause(errors.UserCreateConflict)
	}

	err = uu.userRepo.InsertUser(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	users = append(users, *user)
	return users, nil
}

func (uu *UserUsecase) GetUserByNickname(nickname string) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectUserByNickname(nickname)
	switch {
	case err == pgx.ErrNoRows:
		return nil, errors.Cause(errors.UserNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}
	return user, nil
}

func (uu *UserUsecase) UpdateProfile(user *models.User) *errors.Error {
	err := uu.userRepo.Update(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return errors.Cause(errors.UserProfileConflict)
		}
		return errors.Cause(errors.UserNotExist)
	}

	return nil
}