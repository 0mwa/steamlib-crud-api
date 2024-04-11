package main

import (
	"TestProject/internal"
	"TestProject/internal/entity_handler"
	"TestProject/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	_ "golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {

	logger := internal.NewLogger()
	validate := validator.New(validator.WithRequiredStructEnabled())
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer logger.Sync()

	db := internal.GetBD()
	usersRepoStruct := repository.Users{db}
	sessionsRepoStruct := repository.Sessions{db}
	authStruct := internal.Auth{usersRepoStruct, sessionsRepoStruct, logger}
	http.HandleFunc("/auth", authStruct.Auth)

	entityHandlers := []entity_handler.EntityHandler{
		entity_handler.Games{logger, validate, rdb},
		entity_handler.Publishers{logger, validate, rdb},
		entity_handler.Developers{logger, validate, rdb},
	}
	for _, v := range entityHandlers {
		http.Handle("/"+v.GetPath(), internal.AuthHandler(http.HandlerFunc(v.GetAll), sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/{id}", internal.AuthHandler(http.HandlerFunc(v.Get), sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/add/{id}", internal.AuthHandler(http.HandlerFunc(v.Post), sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/delete/{id}", internal.AuthHandler(http.HandlerFunc(v.Del), sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/update/{id}", internal.AuthHandler(http.HandlerFunc(v.Put), sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/get_counter", internal.AuthHandler(http.HandlerFunc(v.GetCounter), sessionsRepoStruct))
	}

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
