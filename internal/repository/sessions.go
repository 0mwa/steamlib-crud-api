package repository

import (
	"database/sql"
	"errors"
	"time"
)

type Sessions struct {
	db *sql.DB
}

func NewSessions(db *sql.DB) *Sessions {
	return &Sessions{db}
}

func (s Sessions) AddSession(userId int, token string) (err error) {
	expiration := time.Now().Add(48 * time.Hour)
	_, err = s.db.Query("INSERT INTO sessions (user_id, token, expiration) VALUES ($1, $2, $3)", userId, token, expiration)
	if err != nil {
		return err
	}
	return nil
}

func (s Sessions) CheckSessionExp(userId int) (token string, err error) {
	err = s.db.QueryRow("SELECT token FROM sessions WHERE user_id = $1 AND expiration >= NOW()  ORDER BY id DESC LIMIT 1", userId).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("no token")
		}
		return "", err
	}
	return token, nil
}

func (s Sessions) CheckToken(tokenHeader string) (ok bool, err error) {
	var trash string
	err = s.db.QueryRow("SELECT token FROM sessions WHERE token = $1 AND expiration >= NOW()", tokenHeader).Scan(&trash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
