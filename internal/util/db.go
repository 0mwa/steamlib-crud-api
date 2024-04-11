package util

import "database/sql"

func NewDb() *sql.DB {
	file := "postgres://postgres:omw@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", file)
	if err != nil {
		panic(err)
	}
	return db
}
