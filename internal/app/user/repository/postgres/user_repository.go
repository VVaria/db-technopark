package postgres

import (
	"github.com/VVaria/db-technopark/internal/models"
	"github.com/jackc/pgx"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/user"
)

type UserRepository struct {
	conn *pgx.ConnPool
}

func NewUserRepository(conn *pgx.ConnPool) user.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (ur *UserRepository) SelectUsers(user *models.User) ([]models.User, error) {
	var users []models.User
	rows, err := ur.conn.Query(`
			select * 
			from users 
			where Nickname=$1 or Email=$2`,
		user.Nickname, user.Email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.Nickname, &userInfo.FullName, &userInfo.Email, &userInfo.About)
		if err != nil {
			return users, err
		}
		users = append(users, userInfo)
	}

	return users, nil
}

func (ur *UserRepository) InsertUser(user *models.User) error {
	_, err := ur.conn.Exec(`
			INSERT INTO 
			users(nickname, fullname, email, about) 
			VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.FullName, user.Email, user.About)

	return err
}

func (ur *UserRepository) SelectUserByNickname(nickname string) (*models.User, error) {
	query := ur.conn.QueryRow(`
		select * 
		from users 
		where nickname=$1 
		LIMIT 1`, nickname)

	var user models.User
	err := query.Scan(&user.Nickname, &user.FullName, &user.Email, &user.About)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(user *models.User) error {
	query := ur.conn.QueryRow(`
			UPDATE users SET
			fullname=COALESCE(NULLIF($1, ''), fullname),
			email=COALESCE(NULLIF($2, ''), email),
			about=COALESCE(NULLIF($3, ''), about)
			WHERE nickname=$4
			RETURNING nickname, fullname, email, about`,
		user.FullName, user.Email, user.About, user.Nickname)

	err := query.Scan(&user.Nickname, &user.FullName, &user.Email, &user.About)
	if err != nil {
		return err
	}

	return nil
}


func (ur *UserRepository) SelectForumUsers(slug string, params *models.Parameters) ([]*models.User, error) {
	var query *pgx.Rows
	var err error

	if params.Since != "" {
		if params.Desc {
			query, err = ur.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE slug=$1 AND nickname < $2
					ORDER BY nickname DESC
					LIMIT NULLIF($3, 0)`,
				slug, params.Since, params.Limit)
		} else {
			query, err = ur.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE slug=$1 AND nickname > $2
					ORDER BY nickname ASC
					LIMIT NULLIF($3, 0)`,
				slug, params.Since, params.Limit)
		}
	} else {
		if params.Desc {
			query, err = ur.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE slug=$1
					ORDER BY nickname DESC
					LIMIT NULLIF($2, 0)`,
				slug, params.Limit)
		} else {
			query, err = ur.conn.Query(`
					SELECT nickname, fullname, about, email FROM forum_users
					WHERE slug=$1
					ORDER BY nickname ASC
					LIMIT NULLIF($2, 0)`,
				slug, params.Limit)
		}
	}
	var users []*models.User
	if err != nil {
		return users, nil
	}
	defer query.Close()

	for query.Next() {
		var u models.User
		err = query.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
		if err != nil {
			return users, err
		}

		users = append(users, &u)
	}

	return users, err
}