package entity_handler

import "net/http"

const MethodError string = "405 - Wrong method"

type EntityHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Del(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	GetPath() string
}
