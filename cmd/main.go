package main

import (
	"TestProject/internal"
	"TestProject/internal/entity_handler"
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

	entityHandlers := []entity_handler.EntityHandler{
		entity_handler.Games{logger, validate, rdb},
		entity_handler.Publishers{logger, validate, rdb},
		entity_handler.Developers{logger, validate, rdb},
	}

	for _, v := range entityHandlers {
		http.HandleFunc("/"+v.GetPath(), v.GetAll)
		http.HandleFunc("/"+v.GetPath()+"/{id}", v.Get)
		http.HandleFunc("/"+v.GetPath()+"/add/{id}", v.Post)
		http.HandleFunc("/"+v.GetPath()+"/delete/{id}", v.Del)
		http.HandleFunc("/"+v.GetPath()+"/update/{id}", v.Put)
		http.HandleFunc("/"+v.GetPath()+"/get_counter", v.GetCounter)
	}

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
