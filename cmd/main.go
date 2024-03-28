package main

import (
	"TestProject/internal"
	"TestProject/internal/entity_handler"
	"github.com/go-playground/validator/v10"
	_ "golang.org/x/net/html"
	"net/http"
)

func main() {

	logger := internal.NewLogger()
	validate := validator.New(validator.WithRequiredStructEnabled())
	defer logger.Sync() // flushes buffer, if any

	entityHandlers := []entity_handler.EntityHandler{
		entity_handler.Games{logger, validate},
		entity_handler.Publishers{logger, validate},
		entity_handler.Developers{logger, validate},
	}

	for _, v := range entityHandlers {
		http.HandleFunc("/"+v.GetPath(), v.GetAll)
		http.HandleFunc("/"+v.GetPath()+"/{id}", v.Get)
		http.HandleFunc("/"+v.GetPath()+"/add/{id}", v.Post)
		http.HandleFunc("/"+v.GetPath()+"/delete/{id}", v.Del)
		http.HandleFunc("/"+v.GetPath()+"/update/{id}", v.Put)
	}

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
