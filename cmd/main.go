package main

import (
	"TestProject/internal"
	"TestProject/internal/entity_handler"
	"TestProject/internal/repository"
	"TestProject/internal/service"
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
	usersRepoStruct := repository.NewUsers(db)
	sessionsRepoStruct := repository.NewSessions(db)
	errToJson := internal.NewErrToJson(logger)
	authService := service.NewAuthService(usersRepoStruct, sessionsRepoStruct, logger)
	authStruct := internal.Auth{usersRepoStruct, sessionsRepoStruct, authService, logger, errToJson}
	http.HandleFunc("/auth", authStruct.Auth)
	http.HandleFunc("/register", authStruct.CreateUser)

	entityHandlers := []entity_handler.EntityHandler{
		entity_handler.Games{logger, validate, rdb, errToJson},
		entity_handler.Publishers{logger, validate, rdb, errToJson},
		entity_handler.Developers{logger, validate, rdb, errToJson},
	}
	for _, v := range entityHandlers {
		http.Handle("/"+v.GetPath(), internal.AuthHandler(http.HandlerFunc(v.GetAll), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/{id}", internal.AuthHandler(http.HandlerFunc(v.Get), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/add/{id}", internal.AuthHandler(http.HandlerFunc(v.Post), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/delete/{id}", internal.AuthHandler(http.HandlerFunc(v.Del), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/update/{id}", internal.AuthHandler(http.HandlerFunc(v.Put), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/get_counter", internal.AuthHandler(http.HandlerFunc(v.GetCounter), *sessionsRepoStruct))
	}

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
