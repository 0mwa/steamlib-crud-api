package entity_handler

import "net/http"

const GamesCounter string = "GAMES_COUNTER"
const DevelopersCounter string = "DEV_COUNTER"
const PublishersCounter string = "PUB_COUNTER"

type EntityHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Del(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetCounter(w http.ResponseWriter, r *http.Request)
}
