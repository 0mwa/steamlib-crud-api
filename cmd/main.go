package main

import (
	"TestProject/internal/handler"
	entity_handler2 "TestProject/internal/handler/entity_handler"
	"TestProject/internal/middleware"
	"TestProject/internal/repository"
	"TestProject/internal/service"
	"TestProject/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	_ "golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	logger := util.NewLogger()
	validate := validator.New(validator.WithRequiredStructEnabled())
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer logger.Sync()

	db := util.NewDb()
	usersRepoStruct := repository.NewUsers(db)
	sessionsRepoStruct := repository.NewSessions(db)
	errToJson := util.NewErrToJson(logger)
	authService := service.NewAuthService(usersRepoStruct, sessionsRepoStruct, logger)
	authStruct := handler.Auth{logger, usersRepoStruct, sessionsRepoStruct, authService, errToJson}
	http.HandleFunc("/auth", authStruct.Auth)
	http.HandleFunc("/register", authStruct.CreateUser)

	entityHandlers := []entity_handler2.EntityHandler{
		entity_handler2.Games{logger, validate, rdb, errToJson, db},
		entity_handler2.Publishers{logger, validate, rdb, errToJson, db},
		entity_handler2.Developers{logger, validate, rdb, errToJson, db},
	}
	for _, v := range entityHandlers {
		http.Handle("/"+v.GetPath(), middleware.AuthHandler(http.HandlerFunc(v.GetAll), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/{id}", middleware.AuthHandler(http.HandlerFunc(v.Get), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/add/{id}", middleware.AuthHandler(http.HandlerFunc(v.Post), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/delete/{id}", middleware.AuthHandler(http.HandlerFunc(v.Del), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/update/{id}", middleware.AuthHandler(http.HandlerFunc(v.Put), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/get_counter", middleware.AuthHandler(http.HandlerFunc(v.GetCounter), *sessionsRepoStruct))
	}

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
