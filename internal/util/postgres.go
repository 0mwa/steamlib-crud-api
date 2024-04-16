package util

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewPostgres(env *Env) *sql.DB {
	file := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", env.PgUser, env.PgPassword, env.PgHost, env.PgPort)
	db, err := sql.Open("postgres", file)
	if err != nil {
		panic(err)
	}
	return db
}
