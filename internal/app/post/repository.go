package post

import "github.com/VVaria/db-technopark/internal/models"

type PostRepository interface {
	SelectPostById(id int) (*models.Post, error)
	Update(post *models.Post) error
	InsertPosts(posts []*models.Post, threadID int, threadForum string) ([]*models.Post, error)
	SelectThreadPostsFlat(id int, params *models.ThreadPostParameters) ([]*models.Post, error)
	SelectThreadPostsTree(id int, params *models.ThreadPostParameters) ([]*models.Post, error)
	SelectThreadPostsParent(id int, params *models.ThreadPostParameters) ([]*models.Post, error)
}
