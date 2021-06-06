package postgres

import (
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/thread"
)

type ThreadRepository struct {
	conn *pgx.ConnPool
}

func NewThreadRepository(conn *pgx.ConnPool) thread.ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}

func (tr *ThreadRepository) SelectThreadByID(id int) (*models.Thread, error) {
	var thread models.Thread
	query := tr.conn.QueryRow(`
		select id, title, author, forum, message, votes, slug, created 
		from threads
		where id=$1 LIMIT 1;`, id)

	err := query.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (tr *ThreadRepository) InsertThread(thread *models.Thread) error {
	query := tr.conn.QueryRow(`
		INSERT INTO threads (title, author, forum, message, slug, created)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, author, forum, message, votes, slug, created`,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created)

	err := query.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) SelectForumThreads(slug string, params *models.Parameters) ([]*models.Thread, error) {
	var query *pgx.Rows
	var err error

	if params.Since != "" {
		if params.Desc {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1 AND created <= $2
					ORDER BY created DESC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1 AND created >= $2
					ORDER BY created ASC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		}
	} else {
		if params.Desc {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1
					ORDER BY created DESC
					LIMIT $2`,
				slug,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT id, title, author, forum, message, votes, slug, created FROM thread
					WHERE forum=$1
					ORDER BY created ASC
					LIMIT $2`,
				slug,
				params.Limit)
		}
	}
	var threads []*models.Thread
	if err != nil {
		return nil, nil
	}
	defer query.Close()

	for query.Next() {
		var t models.Thread
		err := query.Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &t)
	}

	return threads, nil
}