package entity_handler

import "net/http"

type Developers struct {
}

func (Developers) Get(w http.ResponseWriter, r *http.Request)    {}
func (Developers) GetAll(w http.ResponseWriter, r *http.Request) {}
func (Developers) Post(w http.ResponseWriter, r *http.Request)   {}
func (Developers) Del(w http.ResponseWriter, r *http.Request)    {}
func (Developers) Put(w http.ResponseWriter, r *http.Request)    {}
