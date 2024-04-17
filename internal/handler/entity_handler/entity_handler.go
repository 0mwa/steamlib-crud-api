package entity_handler

import (
	"net/http"
)

const (
	GamesCounter      string = "GAMES_COUNTER"
	DevelopersCounter string = "DEV_COUNTER"
	PublishersCounter string = "PUB_COUNTER"
)

type GameIn struct {
	Name        *string `json:"name" validate:"max=255"`
	Img         *string `json:"img" validate:"max=255"`
	Description *string `json:"description" validate:"max=255"`
	Rating      *int    `json:"rating" validate:"gte=0,lte=10"`
	DeveloperId *int    `json:"developer_id" validate:"gte=0"`
	PublisherId *int    `json:"publisher_id" validate:"gte=0"`
}

type DevPubIn struct {
	Name    *string `json:"name" validate:"max=255"`
	Country *string `json:"country" validate:"max=100"`
}

type Counter struct {
	Count string `json:"count"`
}

type EntityHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Del(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetCounter(w http.ResponseWriter, r *http.Request)
}
