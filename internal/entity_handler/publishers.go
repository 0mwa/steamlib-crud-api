package entity_handler

import "net/http"

type Publishers struct {
}

func (Publishers) Get(w http.ResponseWriter, r *http.Request)    {}
func (Publishers) GetAll(w http.ResponseWriter, r *http.Request) {}
func (Publishers) Post(w http.ResponseWriter, r *http.Request)   {}
func (Publishers) Del(w http.ResponseWriter, r *http.Request)    {}
func (Publishers) Put(w http.ResponseWriter, r *http.Request)    {}
