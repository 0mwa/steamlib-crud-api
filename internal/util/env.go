package util

import "os"

type Env struct {
	CrudApiPort string

	PgUser     string
	PgPassword string
	PgHost     string
	PgPort     string

	RedisUser     string
	RedisPassword string
	RedisHost     string
	RedisPort     string
}

func NewEnv() *Env {
	return &Env{
		os.Getenv("CRUD_API_PORT"),

		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),

		os.Getenv("REDIS_USER"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	}
}
