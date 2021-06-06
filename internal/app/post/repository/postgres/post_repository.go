package postgres

import (
	"fmt"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"
	"strings"
	"time"

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
		Scan(&post.ID, &post.Author, &post.AuthorId, &post.Created, &post.Forum, &post.IsEdited, &post.Message,
		&post.Parent, &post.Thread, &post.Path)
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
		post.Message, post.ID,
	).Scan(&post.ID, &post.Author, &post.AuthorId, &post.Created, &post.Forum, &post.IsEdited, &post.Message,
		&post.Parent, &post.Thread, &post.Path)
	return err
}

func (pr *PostRepository) InsertPosts(posts []*models.Post, threadID int, threadForum string) ([]*models.Post, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	queryString := `INSERT INTO posts(author, created, forum, message, parent, id_thread) VALUES `
	var queryParameters []interface{}
	created := time.Now()
	for i, p := range posts {
		value := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d),",
			i * 6 + 1, i * 6 + 2, i * 6 + 3, i * 6 + 4, i * 6 + 5, i * 6 + 6,
		)

		queryString += value
		queryParameters = append(queryParameters, p.Author, created, threadForum, p.Message, p.Parent, threadID)
	}

	queryString = strings.TrimSuffix(queryString, ",")
	queryString += ` RETURNING *`

	query, err := pr.conn.Query(queryString, queryParameters...)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	postsInfo := make([]*models.Post, 0)
	for query.Next() {
		var post models.Post

		err := query.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.Message, &post.IsEdited,
			&post.Parent, &post.Thread, &post.Path)
		if err != nil {
			return nil, err
		}

		postsInfo = append(postsInfo, &post)
	}

	return postsInfo, nil
}