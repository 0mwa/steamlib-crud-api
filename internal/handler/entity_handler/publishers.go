package entity_handler

import (
	"Crud-Api/internal/entity"
	"Crud-Api/internal/handler"
	"Crud-Api/internal/repository"
	"Crud-Api/internal/util"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Publishers struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
	Rds       *redis.Client
	ErrTo     *util.ErrToJson
	Db        *sql.DB
}

func (p Publishers) GetPath() string {
	return "publishers"
}

func (p Publishers) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}
	var err error
	var result *sql.Rows
	var marshaled []byte

	id := r.PathValue("id")
	result, err = p.Db.Query(repository.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	pub := entity.Publisher{}
	result.Next()
	err = result.Scan(&pub.Name, &pub.Country, &pub.SteamId)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	marshaled, err = json.Marshal(pub)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
}
func (p Publishers) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	result, err = p.Db.Query(repository.SelectPublishers)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	pubs := make([]entity.Publisher, 0)
	pub := entity.Publisher{}
	for result.Next() {
		err = result.Scan(&pub.Name, &pub.Country, &pub.SteamId)
		if err != nil {
			p.Logger.Error(err)
			p.ErrTo.ErrToJson(w, err)
			return
		}
		pubs = append(pubs, pub)
	}
	marshaled, err = json.Marshal(pubs)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
}
func (p Publishers) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	id := r.PathValue("id")
	var response *SteamResponse
	response = &SteamResponse{GameList: make(map[string]SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	pubNameArr := response.GameList[id].Data.Publishers
	pubCountry := randomdata.Country(randomdata.FullCountry)
	if len(pubNameArr) > 0 {
		pubName := strings.Join(pubNameArr, " ")
		_, err = p.Db.Query("INSERT INTO publishers (name, country, steam_id) VALUES ($1, $2, $3)", pubName, pubCountry, id)
		if err != nil {
			if strings.Contains(err.Error(), "\"pubsteam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				p.Logger.Error(errors.New("409 - Publisher already exists"))
				p.ErrTo.ErrToJson(w, errors.New("409 - Publisher already exists"))
			} else {
				p.Logger.Error(err)
				p.ErrTo.ErrToJson(w, err)
				return
			}
		}
		p.Rds.Del(context.Background(), PublishersCounter)
		p.Logger.Infof("Pub with id: %s added. %s is flushed", id, PublishersCounter)
	} else {
		w.WriteHeader(http.StatusConflict)
		p.Logger.Error(errors.New("409 - No publisher with such id"))
		p.ErrTo.ErrToJson(w, errors.New("409 - No publisher with such id"))
	}
}
func (p Publishers) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var err error
	var response *sql.Rows

	id := r.PathValue("id")
	response, err = p.Db.Query(repository.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		p.Logger.Error(errors.New("409 - no publisher to delete with such id"))
		p.ErrTo.ErrToJson(w, errors.New("409 - no publisher to delete with such id"))
		return
	}
	_, err = p.Db.Query(repository.DeletePublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	p.Rds.Del(context.Background(), PublishersCounter)
	p.Logger.Infof("Pub with id: %s deleted. %s is flushed", id, PublishersCounter)
}
func (p Publishers) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	pubStruct := DevPubIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(requestBody, &pubStruct)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}

	if err = p.Validator.Struct(&pubStruct); err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}

	id := r.PathValue("id")
	response, err = p.Db.Query(repository.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		p.Logger.Error(errors.New("409 - no publisher to update with such id"))
		p.ErrTo.ErrToJson(w, errors.New("409 - no publisher to update with such id"))
		return
	}
	_, err = p.Db.Query(repository.UpdatePublisherById, pubStruct.Name, pubStruct.Country, id)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
}

func (p Publishers) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}
	count := Counter{}
	var marshaled []byte
	var err error
	count.Count, err = p.Rds.Get(context.Background(), PublishersCounter).Result()
	if errors.Is(err, redis.Nil) {
		err = p.Db.QueryRow(repository.GetPublishersCount).Scan(&count.Count)
		p.Logger.Infof("Redis key %s is empty, getting data from DB", PublishersCounter)
		if err != nil {
			p.Logger.Error(err)
			p.ErrTo.ErrToJson(w, err)
			return
		}
		p.Rds.Append(context.Background(), PublishersCounter, count.Count)
		marshaled, err = json.Marshal(count)
		if err != nil {
			p.Logger.Error(err)
			p.ErrTo.ErrToJson(w, err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			p.Logger.Error(err)
			p.ErrTo.ErrToJson(w, err)
			return
		}
		return
	}

	p.Logger.Infof("Redis key %s is found, getting data from Redis", PublishersCounter)
	marshaled, err = json.Marshal(count)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		p.ErrTo.ErrToJson(w, err)
		return
	}
	return
}
