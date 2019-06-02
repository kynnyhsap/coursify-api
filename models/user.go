package models

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"user_name"`
	FullName    string    `json:"full_name"`
	Avatar      string    `json:"avatar"`
	About       string    `json:"about"`
	DateCreated time.Time `json:"date_created"`
}

type Credentials struct {
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

type UserCreateInput struct {
	About    string `json:"about"`
	Avatar   string `json:"avatar"`
	FullName string `json:"full_name"`
	Credentials
}

type ModelUser struct {
	model
}

type IUserGetter interface {
	Get(id int64) User
}

type IUserLister interface {
	GetList(limit, offset int, search string) []User
	Count() int
}

type IUserCreator interface {
	Create(in UserCreateInput) int64
	IUserGetter
}

func (m ModelUser) Get(id int64) User {
	user := User{}

	row := m.db.QueryRow(`
		SELECT
			id,
			user_name,
		    full_name,
		    avatar,
		    about,
		    date_created
		FROM users WHERE id = ?
	`, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.FullName,
		&user.Avatar,
		&user.About,
		&user.DateCreated)

	if err != nil {
		return User{}
	}

	return user
}

func (m ModelUser) GetList(limit, offset int, search string) []User {
	users := make([]User, 0)

	rows, err := m.db.Query(`
		SELECT
		       id, user_name, full_name, avatar, about
		FROM users
		WHERE full_name LIKE ?
		LIMIT ? OFFSET ?
	`, "%"+search+"%", limit, offset)

	if err != nil {
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.Scan(&user.ID, &user.Name, &user.FullName, &user.Avatar, &user.About)
		if err != nil {
			//
			return users
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return users
}

func (m ModelUser) Create(in UserCreateInput) int64 {
	stmt, err := m.db.Prepare(`
		INSERT INTO users (
			full_name,
			user_name,
			avatar,
		    about,
			password_hash,
			date_created
		) VALUE (?, ?, ?, ?, ?, NOW())
	`)

	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.FullName, in.UserName, in.Avatar, in.About, in.PasswordHash)
	if err != nil {
		return 0
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	return lastID
}

func (m ModelUser) Count() int {
	rows, err := m.db.Query(`SELECT COUNT(*) FROM users`)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Println(err)
			return 0
		}
	}

	return count
}

func NewUserModel(db *sql.DB) ModelUser {
	return ModelUser{model{db}}
}
