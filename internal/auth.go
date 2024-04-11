package internal

import (
	"TestProject/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

type Auth struct {
	UsersRepo repository.Users
	SessRepo  repository.Sessions
	Logger    *zap.SugaredLogger
}

type reqBody struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}

func (a Auth) errToJson(w http.ResponseWriter, externalError error) {
	errrr := ErrOut{externalError.Error()}
	marshaled, err := json.Marshal(errrr)
	if err != nil {
		a.Logger.Error(err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		a.Logger.Error(err)
		return
	}
}

func (a Auth) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		a.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var requestBody []byte
	req := reqBody{}
	requestBody, err = io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		a.Logger.Error(err)
		a.errToJson(w, err)
		return
	}
	userId, err := a.UsersRepo.CheckUser(req.Login, req.Passwd)
	if err != nil {
		a.Logger.Error(err)
		a.errToJson(w, err)
		return
	}
	if userId != 0 {
		var token string
		token, err = a.SessRepo.CheckSessionExp(userId)
		if err != nil {
			errstr := err.Error()
			if errstr == "no token" {
				token = strconv.Itoa(userId+rand.Int()) + randomdata.FirstName(randomdata.RandomGender)
				err = a.SessRepo.AddSession(userId, token)
				if err != nil {
					a.Logger.Error(err)
					a.errToJson(w, err)
					return
				}
				_, err = w.Write([]byte(fmt.Sprintf(`{"msg":"New token given.","token":"%s"}`, token)))
				if err != nil {
					a.Logger.Error(err)
					a.errToJson(w, err)
					return
				}
				a.Logger.Infof("New token given to user with id %v. Case %v.", userId, errstr)
				return
			}
			a.Logger.Error(err)
			a.errToJson(w, err)
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf(`{"msg":"Token up to date.","token":"%s"}`, token)))
		if err != nil {
			a.Logger.Error(err)
			a.errToJson(w, err)
			return
		}
		a.Logger.Infof("%v User's token up to date.", userId)
		return
	}
	a.Logger.Info("Authentication failed.")
	a.errToJson(w, errors.New("Authentication failed."))
}
