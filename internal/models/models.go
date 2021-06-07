package models

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"time"
)

type Forum struct {
	ID      int    `json:"-"`
	Slug    string `json:"slug"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type Thread struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Created  time.Time `json:"created"`
	Author   string    `json:"author"`
	Forum    string    `json:"forum"`
	Message  string    `json:"message"`
	Slug     string    `json:"slug"`
	Votes    int       `json:"votes"`
}

type User struct {
	ID       int    `json:"-"`
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Post struct {
	ID       int              `json:"id"`
	Author   string           `json:"author"`
	AuthorId int              `json:"-"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	Message  string           `json:"message"`
	Parent   JsonNullInt64    `json:"parent"`
	Thread   int              `json:"thread"`
	Path     pgtype.Int8Array `json:"-"`
	IsEdited bool             `json:"isEdited"`
}

type PostShort struct {
	ID       int              `json:"id"`
	Author   string           `json:"author"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	Message  string           `json:"message"`
	Thread   int              `json:"thread"`
	Parent   JsonNullInt64    `json:"parent"`
}

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"-"`
}

type AllPostInfo struct {
	Author *User       `json:"author"`
	Forum  *Forum      `json:"forum"`
	Post   *Post       `json:"post"`
	Thread *Thread 	   `json:"thread"`
}

type Parameters struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type ThreadPostParameters struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Desc  bool   `json:"desc"`
	Sort  string `json:"sort"`
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type AllPostInfoWithoutSlug struct {
	Author *User       `json:"author"`
	Forum  *Forum      `json:"forum"`
	Post   *Post       `json:"post"`
	Thread *ThreadWithoutSlug 	   `json:"thread"`
}

type ThreadWithoutSlug struct {
	Id      int    `json:"id"`
	Author  string `json:"author"`
	Created time.Time `json:"created"`
	Forum   string `json:"forum"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Votes   int    `json:"votes"`
}

func ThreadNoSlug(thread *Thread) *ThreadWithoutSlug {
	return &ThreadWithoutSlug{
		Id:      thread.Id,
		Author:  thread.Author,
		Created: thread.Created,
		Forum:   thread.Forum,
		Title:   thread.Title,
		Message: thread.Message,
		Votes:   thread.Votes,
	}
}


func IsUuid(value string) bool {
	n := len(value)

	if n > 36 || n < 32 {
		return false
	}

	_, err := uuid.Parse(value)

	return err == nil
}