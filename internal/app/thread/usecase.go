package thread

import (
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, *errors.Error)
}
