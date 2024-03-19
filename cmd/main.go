package main

import (
	"TestProject/internal"
	"TestProject/internal/entity_handler"
	_ "golang.org/x/net/html"
	"net/http"
)

func main() {

	logger := internal.NewLogger()
	defer logger.Sync() // flushes buffer, if any

	a := entity_handler.Games{logger}

	http.HandleFunc("/games", a.GetAll)
	http.HandleFunc("/games/{id}", a.Get)
	http.HandleFunc("/games/add/{id}", a.Post)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
