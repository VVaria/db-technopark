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
	err := pr.conn.QueryRow(`SELECT * FROM posts WHERE id=$1 LIMIT 1`, id).
		Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.Message,
		&post.Parent, &post.Thread, &post.Path, &post.IsEdited)
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
	).Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.Message,
		&post.Parent, &post.Thread, &post.Path, &post.IsEdited)
	return err
}

func (pr *PostRepository) InsertPosts(posts []*models.Post, threadID int, threadForum string) ([]*models.Post, error) {
	if len(posts) == 0 {
		return posts, nil
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

		err := query.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.Message, &post.Parent,
			&post.Thread, &post.Path, &post.IsEdited)
		if err != nil {
			return nil, err
		}

		postsInfo = append(postsInfo, &post)
	}

	if pgErr, ok := query.Err().(pgx.PgError); ok {
		if pgErr.Code == "00409" {
			return nil, query.Err()
		}

		if pgErr.Code == "23503" {
			return nil, query.Err()
		}
	}
	return postsInfo, nil
}

func (pr *PostRepository) SelectThreadPostsFlat(id int, params *models.ThreadPostParameters) ([]*models.Post, error) {
	var query *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread=$1 
					ORDER BY id DESC 
					LIMIT NULLIF($2, 0)`,
				id, params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread=$1 
					ORDER BY id ASC 
					LIMIT NULLIF($2, 0)`,
				id, params.Limit)
		}
	} else {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread=$1 AND id < $2 
					ORDER BY id DESC 
					LIMIT NULLIF($3, 0)`,
				id, params.Since, params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts 
					WHERE id_thread=$1 AND id > $2 
					ORDER BY id ASC 
					LIMIT NULLIF($3, 0)`,
				id, params.Since, params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var posts []*models.Post
	for query.Next() {
		var p models.Post
		err = query.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.Message, &p.Parent, &p.Thread, &p.Path, &p.IsEdited)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, err
}

func (pr *PostRepository) SelectThreadPostsTree(id int, params *models.ThreadPostParameters) ([]*models.Post, error) {
	var query *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread = $1 
					ORDER BY path DESC, id DESC 
					LIMIT $2`,
				id,
				params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread = $1 
					ORDER BY path ASC, id ASC
					LIMIT $2`,
				id,
				params.Limit)
		}
	} else {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread = $1 AND PATH < (SELECT path FROM posts WHERE id = $2)
					ORDER BY path DESC, id DESC
					LIMIT $3`,
				id,
				params.Since,
				params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE id_thread = $1 AND PATH > (SELECT path FROM posts WHERE id = $2)
					ORDER BY path ASC, id ASC
					LIMIT $3`,
				id,
				params.Since,
				params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var posts []*models.Post
	for query.Next() {
		var p models.Post
		err = query.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.Message, &p.Parent, &p.Thread, &p.Path, &p.IsEdited)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

func (pr *PostRepository) SelectThreadPostsParent(id int, params *models.ThreadPostParameters) ([]*models.Post, error) {
	var query *pgx.Rows
	var err error

	if params.Since == 0 {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE id_thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path, id`,
				id, params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE id_thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
					ORDER BY path, id`,
				id, params.Limit)
		}
	} else {
		if params.Desc {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE id_thread = $1 AND parent IS NULL AND PATH[1] <
					(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) 
					ORDER BY path[1] DESC, path, id`,
				id, params.Since, params.Limit)
		} else {
			query, err = pr.conn.Query(`
					SELECT * FROM posts
					WHERE path[1] IN (SELECT id FROM posts WHERE id_thread = $1 AND parent IS NULL AND PATH[1] >
					(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path, id`,
				id, params.Since, params.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var posts []*models.Post
	for query.Next() {
		var p models.Post
		err = query.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.Message, &p.Parent, &p.Thread, &p.Path, &p.IsEdited)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}