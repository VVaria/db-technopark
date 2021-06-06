package thread

import (
	"github.com/VVaria/db-technopark/internal/models"
)

type ThreadRepository interface {
	SelectThreadByID(id int) (*models.Thread, error)
	InsertThread(thread *models.Thread) error
	SelectForumThreads(slug string, params *models.Parameters) ([]*models.Thread, error)
}
