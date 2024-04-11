package repository

import (
	"database/sql"
	"errors"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (u Users) GetUserId(login string, passwd string) (userId int, err error) {
	var result *sql.Rows
	result, err = u.db.Query("SELECT id FROM users WHERE login = $1 AND passwd = $2", login, passwd)
	if err != nil {
		return 0, err
	}
	if result.Next() {
		err = result.Scan(&userId)
		if err != nil {
			return 0, err
		}
		return userId, nil
	}
	return 0, nil
}

func (u Users) CreateUser(login string, passwd string) (err error) {
	_, err = u.db.Query("INSERT INTO users (login, passwd) VALUES ($1, $2)", login, passwd)
	if err != nil {
		return err
	}
	return nil
}

func (u Users) CheckUserByLogin(login string) (ok bool, err error) {
	var trash string
	err = u.db.QueryRow("SELECT login FROM users WHERE login = $1", login).Scan(&trash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
