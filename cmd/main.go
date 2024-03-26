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

	a := entity_handler.Games{logger, validate}
	b := entity_handler.Publishers{logger, validate}

	http.HandleFunc("/games", a.GetAll)
	http.HandleFunc("/games/{id}", a.Get)
	http.HandleFunc("/games/add/{id}", a.Post)
	http.HandleFunc("/games/delete/{id}", a.Del)
	http.HandleFunc("/games/update/{id}", a.Put)

	http.HandleFunc("/publishers", b.GetAll)
	http.HandleFunc("/publishers/{id}", b.Get)
	http.HandleFunc("/publishers/add/{id}", b.Post)
	http.HandleFunc("/publishers/delete/{id}", b.Del)
	http.HandleFunc("/publishers/update/{id}", b.Put)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
