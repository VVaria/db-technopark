package postgres

import (
	"fmt"

	"github.com/VVaria/db-technopark/internal/app/thread"
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"
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
		select id, title, created, author, forum, message, slug, votes
		from threads
		where id=$1 LIMIT 1`, id)

	err := query.Scan(&thread.Id, &thread.Title, &thread.Created, &thread.Author, &thread.Forum, &thread.Message,
		&thread.Slug, &thread.Votes)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (tr *ThreadRepository) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread :=  &models.Thread{}
	query := tr.conn.QueryRow(`
		select id, title, created, author, forum, message, slug, votes
		from threads
		where slug=$1 LIMIT 1`, slug)

	err := query.Scan(&thread.Id, &thread.Title, &thread.Created, &thread.Author, &thread.Forum, &thread.Message,
		&thread.Slug, &thread.Votes)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (tr *ThreadRepository) InsertThread(thread *models.Thread) error {
	query := tr.conn.QueryRow(`
		INSERT INTO threads (title, created, author, forum, message, slug)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, created, author, forum, message, slug, votes`,
		thread.Title,
		thread.Created,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug)

	err := query.Scan(&thread.Id, &thread.Title, &thread.Created, &thread.Author, &thread.Forum, &thread.Message,
		&thread.Slug, &thread.Votes)
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
					SELECT * FROM threads
					WHERE forum=$1 AND created <= $2
					ORDER BY created DESC
					LIMIT $3`,
				slug,
				params.Since,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT * FROM threads
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
					SELECT * FROM threads
					WHERE forum=$1
					ORDER BY created DESC
					LIMIT $2`,
				slug,
				params.Limit)
		} else {
			query, err = tr.conn.Query(`
					SELECT * FROM threads
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
		err := query.Scan(&t.Id, &t.Title, &t.Created, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &t)
	}

	return threads, nil
}

func (tr *ThreadRepository) UpdateThread(thread *models.Thread) error {
	queryString := `
		UPDATE threads 
		SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) 
		WHERE %s 
		RETURNING * `

	var query *pgx.Row
	if thread.Slug == "" {
		queryString = fmt.Sprintf(queryString, `id=$3`)
		query = tr.conn.QueryRow(queryString, thread.Title, thread.Message, thread.Id)
	} else {
		queryString = fmt.Sprintf(queryString, `slug=$3`)
		query = tr.conn.QueryRow(queryString, thread.Title, thread.Message, thread.Slug)
	}

	err := query.Scan(&thread.Id, &thread.Title, &thread.Created, &thread.Author, &thread.Forum, &thread.Message, &thread.Slug,
		 &thread.Votes)

	if err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) InsertVote(vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			INSERT INTO votes(author, voice, id_thread) 
			VALUES ($1, $2, $3)`,
		vote.Nickname, vote.Voice, vote.Thread)
	if err != nil {
		return err
	}
	//_ = tr.UpdateThreadVotes(vote.Thread)
	return nil
}

func (tr *ThreadRepository) UpdateVote(vote *models.Vote) error {
	_, err := tr.conn.Exec(`
			UPDATE votes 
			SET voice=$1 
			WHERE author=$2 and id_thread=$3`,
		vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return err
	}
	//_ = tr.UpdateThreadVotes(vote.Thread)
	return nil
}

func (tr *ThreadRepository) UpdateThreadVotes(id int) error {
	_, err := tr.conn.Exec(`
			UPDATE threads 
			SET voice=voice+1 
			WHERE id=$3`,
		id)
	if err != nil {
		return err
	}
	return nil
}