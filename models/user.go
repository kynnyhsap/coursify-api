package models

import "database/sql"

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Credentials struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

type UserCreateInput struct {
	Avatar string `json:"avatar"`

	Credentials
}

type ModelUser struct {
	model
}

type IUserGetter interface {
	Get(id int64) User
	GetByLogin(login string) User
}

type IUserCreator interface {
	Create(in UserCreateInput) int64
	IUserGetter
}

func (m ModelUser) Get(id int64) User {
	user := User{}

	row := m.db.QueryRow(`SELECT id, name, avatar FROM users WHERE id = ?`, id)
	err := row.Scan(&user.ID, &user.Name, &user.Avatar)
	if err != nil {
		return User{}
	}

	return user
}

func (m ModelUser) GetByLogin(login string) User {
	user := User{}

	row := m.db.QueryRow(`SELECT id, name, avatar FROM users WHERE name = ?`, login)
	err := row.Scan(&user.ID, &user.Name, &user.Avatar)
	if err != nil {
		return User{}
	}

	return user
}

func (m ModelUser) Create(in UserCreateInput) int64 {
	stmt, err := m.db.Prepare(`INSERT INTO users(name, avatar, password_hash) VALUE(?, ?, ?)`)
	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.Name, in.Avatar, in.PasswordHash)
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
