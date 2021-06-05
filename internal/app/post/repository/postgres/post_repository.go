package postgres

import (
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"
	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/post"
	//"github.com/VVaria/db-technopark/internal/app/tools/null"
)

type PostRepository struct {
	conn *pgx.ConnPool
}

func NewPostRepository(conn *pgx.ConnPool) post.PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

func (pr *PostRepository) SelectPostById(id int) (*models.Post, error) {
	var post models.Post

	err := pr.conn.QueryRow(`SELECT * FROM posts WHERE id=$1 LIMIT 1;`, id).
		Scan(
		&post.ID,
		&post.Author,
		&post.AuthorId,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&post.Parent,
		&post.Thread,
		&post.Path)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pr *PostRepository) Update(post *models.Post) error {
	err := pr.conn.QueryRow(
		`UPDATE posts 
			SET message=COALESCE(NULLIF($1, ''), message),
		 	isEdited = CASE WHEN $1 = '' OR message = $1 THEN isEdited ELSE true END
			WHERE id=$2 
			RETURNING *`,
		post.Message,
		post.ID,
	).Scan(
		&post.ID,
		&post.Author,
		&post.AuthorId,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&post.Parent,
		&post.Thread,
		&post.Path,
	)
	return err
}