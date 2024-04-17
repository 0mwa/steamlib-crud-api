package handler

import (
	"Crud-Api/internal/repository"
	"Crud-Api/internal/service"
	"Crud-Api/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Auth struct {
	Logger      *zap.SugaredLogger
	UsersRepo   *repository.Users
	SessRepo    *repository.Sessions
	AuthService *service.Auth
	ErrTo       *util.ErrToJson
}

type reqBody struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}

func (a Auth) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		a.ErrTo.ErrToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var requestBody []byte
	var token string
	req := reqBody{}
	requestBody, err = io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		a.Logger.Error(err)
		a.ErrTo.ErrToJson(w, err)
		return
	}

	token, err = a.AuthService.Auth(req.Login, req.Passwd)
	if err != nil {
		a.Logger.Error(err)
		a.ErrTo.ErrToJson(w, err)
		return
	}
	if token == "" {
		a.Logger.Warn("Authentication failed")
		a.ErrTo.ErrToJson(w, errors.New("Authentication failed"))
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf(`{"msg":"Success","token":"%s"}`, token)))
	if err != nil {
		a.Logger.Error(err)
		a.ErrTo.ErrToJson(w, err)
		return
	}
	a.Logger.Infof("Token given to user")
}

func (a Auth) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		a.ErrTo.ErrToJson(w, errors.New(MethodError))
		return
	}
	var err error
	var requestBody []byte
	req := reqBody{}
	requestBody, err = io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		a.Logger.Error(err)
		a.ErrTo.ErrToJson(w, err)
		return
	}
	err = a.AuthService.CreateUser(req.Login, req.Passwd)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			a.Logger.Warn(err)
			_, err = w.Write([]byte(`{"msg":"User with this login already exists"}`))
			return
		}
		a.Logger.Error(err)
		a.ErrTo.ErrToJson(w, err)
		return
	}
	a.Logger.Infof("User created with login: %v", req.Login)
	_, err = w.Write([]byte(`{"msg":"Registration success"}`))
}
