package entity_handler

import (
	"TestProject/internal"
	"TestProject/internal/entity"
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

type Developers struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
	Rds       *redis.Client
	ErrTo     *internal.ErrToJson
}

func (d Developers) GetPath() string {
	return "developers"
}

func (d Developers) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}
	var err error
	var result *sql.Rows
	var marshaled []byte

	id := r.PathValue("id")
	db := internal.GetBD()
	result, err = db.Query(internal.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	dev := entity.Developer{}
	result.Next()
	err = result.Scan(&dev.Name, &dev.Country, &dev.SteamId)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	marshaled, err = json.Marshal(dev)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
}
func (d Developers) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	db := internal.GetBD()
	result, err = db.Query(internal.SelectDevelopers)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	devs := make([]entity.Developer, 0)
	dev := entity.Developer{}
	for result.Next() {
		err = result.Scan(&dev.Name, &dev.Country, &dev.SteamId)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, err)
			return
		}
		devs = append(devs, dev)
	}
	marshaled, err = json.Marshal(devs)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
}
func (d Developers) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	var response *internal.SteamResponse
	response = &internal.SteamResponse{GameList: make(map[string]internal.SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	devNameArr := response.GameList[id].Data.Publishers
	devCountry := randomdata.Country(randomdata.FullCountry)
	if len(devNameArr) > 0 {
		devName := strings.Join(devNameArr, " ")
		_, err = db.Query("INSERT INTO developers (name, country, steam_id) VALUES ($1, $2, $3)", devName, devCountry, id)
		if err != nil {
			if strings.Contains(err.Error(), "\"devsteam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				d.Logger.Error(errors.New("409 - Developer already exists"))
				d.ErrTo.ErrToJson(w, errors.New("409 - Developer already exists"))
			} else {
				d.Logger.Error(err)
				d.ErrTo.ErrToJson(w, err)
				return
			}
		}
		d.Rds.Del(context.Background(), DevelopersCounter)
		d.Logger.Infof("Dev with id: %s added. %s is flushed", id, DevelopersCounter)
	} else {
		w.WriteHeader(http.StatusConflict)
		d.Logger.Error(errors.New("409 - No developer with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - No developer with such id"))
	}
}
func (d Developers) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}

	var err error
	var response *sql.Rows

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		d.Logger.Error(errors.New("409 - no developer to delete with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - no developer to delete with such id"))
		return
	}
	_, err = db.Query(internal.DeleteDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	d.Rds.Del(context.Background(), DevelopersCounter)
	d.Logger.Infof("Dev with id: %s deleted. %s is flushed", id, DevelopersCounter)
}
func (d Developers) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	devStruct := internal.DevPubIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(requestBody, &devStruct)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}

	if err = d.Validator.Struct(&devStruct); err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		d.Logger.Error(errors.New("409 - no developer to update with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - no developer to update with such id"))
		return
	}
	_, err = db.Query(internal.UpdatePublisherById, devStruct.Name, devStruct.Country, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
}

func (d Developers) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(internal.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, errors.New(internal.MethodError))
		return
	}
	count := internal.Counter{}
	var marshaled []byte
	var err error
	count.Count, err = d.Rds.Get(context.Background(), DevelopersCounter).Result()
	if errors.Is(err, redis.Nil) {
		db := internal.GetBD()
		err = db.QueryRow(internal.GetDevelopersCount).Scan(&count.Count)
		d.Logger.Infof("Redis key %s is empty, getting data from DB", DevelopersCounter)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, err)
			return
		}
		d.Rds.Append(context.Background(), DevelopersCounter, count.Count)
		marshaled, err = json.Marshal(count)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, err)
			return
		}
		return
	}

	d.Logger.Infof("Redis key %s is found, getting data from Redis", DevelopersCounter)
	marshaled, err = json.Marshal(count)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, err)
		return
	}
	return
}
