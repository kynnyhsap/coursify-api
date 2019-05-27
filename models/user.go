package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"user_name"`
	FullName    string    `json:"full_name"`
	Avatar      string    `json:"avatar"`
	DateCreated time.Time `json:"date_created"`
}

type Credentials struct {
	Name         string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

type UserCreateInput struct {
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
		    date_created
		FROM users WHERE id = ?
	`, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.FullName,
		&user.Avatar,
		&user.DateCreated)

	if err != nil {
		return User{}
	}

	return user
}

func (m ModelUser) Create(in UserCreateInput) int64 {
	stmt, err := m.db.Prepare(`
		INSERT INTO users (
			full_name,
			user_name,
			avatar,
			password_hash,
			date_created
		) VALUE (?, ?, ?, ?, NOW())
	`)

	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.FullName, in.Name, in.Avatar, in.PasswordHash)
	if err != nil {
		return 0
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	return lastID
}

func NewUserModel(db *sql.DB) ModelUser {
	return ModelUser{model{db}}
}
