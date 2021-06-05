package postgres

import (
	"database/sql"
	"github.com/VVaria/db-technopark/internal/models"

	//"github.com/VVaria/db-technopark/internal/app/models"
	"github.com/VVaria/db-technopark/internal/app/user"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (ur *UserRepository) SelectUsers(user *models.User) ([]models.User, error) {
	var users []models.User
	rows, err := ur.conn.Query(`
			select Nickname, FullName, About, Email 
			from users 
			where Nickname=$1 or Email=$2;
			`, user.Nickname, user.Email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.Nickname, &userInfo.FullName, &userInfo.About, &userInfo.Email)
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
			users(nickname, fullname, about, email) 
			VALUES ($1, $2, $3, $4)`,
		user.Nickname,
		user.FullName,
		user.About,
		user.Email)

	return err
}

func (ur *UserRepository) SelectUserByNickname(nickname string) (*models.User, error) {
	query := ur.conn.QueryRow(`
		select nickname, fullname, about, email 
		from users 
		where nickname=$1 
		LIMIT 1;
		`, nickname)

	var user models.User
	err := query.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(user *models.User) error {
	query := ur.conn.QueryRow(`
			UPDATE users SET
			fullname=COALESCE(NULLIF($1, ''), fullname),
			about=COALESCE(NULLIF($2, ''), about),
			email=COALESCE(NULLIF($3, ''), email)
			WHERE nickname=$4
			RETURNING nickname, fullname, about, email
			`, user.FullName, user.About, user.Email, user.Nickname)

	err := query.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
	if err != nil {
		return err
	}

	return nil
}