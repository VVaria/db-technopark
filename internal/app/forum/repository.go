package forum

import "github.com/VVaria/db-technopark/internal/models"

type ForumRepository interface {
	SelectForumBySlug(slug string) (*models.Forum, error)
	CreateForum(forum *models.Forum) error
}
