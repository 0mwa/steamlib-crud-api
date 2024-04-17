package main

import (
	"Crud-Api/internal/handler"
	"Crud-Api/internal/handler/entity_handler"
	"Crud-Api/internal/middleware"
	"Crud-Api/internal/repository"
	"Crud-Api/internal/service"
	"Crud-Api/internal/util"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func main() {
	logger := util.NewLogger()
	defer logger.Sync()
	validate := validator.New(validator.WithRequiredStructEnabled())
	env := util.NewEnv()
	rdb := util.NewRedis(env)
	db := util.NewPostgres(env)
	usersRepoStruct := repository.NewUsers(db)
	sessionsRepoStruct := repository.NewSessions(db)
	errToJson := util.NewErrToJson(logger)
	authService := service.NewAuthService(usersRepoStruct, sessionsRepoStruct, logger)

	authHandler := handler.Auth{logger, usersRepoStruct, sessionsRepoStruct, authService, errToJson}

	http.HandleFunc("/auth", authHandler.Auth)
	http.HandleFunc("/register", authHandler.CreateUser)

	entityHandlers := []entity_handler.EntityHandler{
		entity_handler.Games{logger, validate, rdb, errToJson, db},
		entity_handler.Publishers{logger, validate, rdb, errToJson, db},
		entity_handler.Developers{logger, validate, rdb, errToJson, db},
	}
	for _, v := range entityHandlers {
		http.Handle("/"+v.GetPath(), middleware.AuthHandler(http.HandlerFunc(v.GetAll), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/{id}", middleware.AuthHandler(http.HandlerFunc(v.Get), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/add/{id}", middleware.AuthHandler(http.HandlerFunc(v.Post), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/delete/{id}", middleware.AuthHandler(http.HandlerFunc(v.Del), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/update/{id}", middleware.AuthHandler(http.HandlerFunc(v.Put), *sessionsRepoStruct))
		http.Handle("/"+v.GetPath()+"/get_counter", middleware.AuthHandler(http.HandlerFunc(v.GetCounter), *sessionsRepoStruct))
	}

	err := http.ListenAndServe(":"+env.CrudApiPort, nil)
	if err != nil {
		panic(err)
	}
}
