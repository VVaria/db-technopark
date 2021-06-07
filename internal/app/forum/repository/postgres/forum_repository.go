package postgres

import (
	"github.com/VVaria/db-technopark/internal/app/forum"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"
	//"github.com/VVaria/db-technopark/internal/app/models"
)

type ForumRepository struct {
	conn *pgx.ConnPool
}

func NewForumRepository(conn *pgx.ConnPool) forum.ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}

func (fr *ForumRepository) SelectForumBySlug(slug string) (*models.Forum, error) {
	var forum models.Forum
	query := fr.conn.QueryRow(`
		select slug, "user", title, posts, threads 
		from forum
		where slug=$1 LIMIT 1`, slug)
	err := query.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (fr *ForumRepository) CreateForum(forum *models.Forum) error {
	query := fr.conn.QueryRow(`
			INSERT INTO forum(title, "user", slug) 
			VALUES ($1, $2, $3)
			RETURNING title, "user", slug, posts, threads`,
		forum.Title, forum.User, forum.Slug)

	err := query.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return err
	}
	return nil
}