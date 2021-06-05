package post

import "github.com/VVaria/db-technopark/internal/models"

type PostRepository interface {
	SelectPostById(id int) (*models.Post, error)
	Update(post *models.Post) error
}
