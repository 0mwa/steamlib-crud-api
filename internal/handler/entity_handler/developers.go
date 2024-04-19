package entity_handler

import (
	"Crud-Api/internal/entity"
	"Crud-Api/internal/repository"
	"Crud-Api/internal/util"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Developers struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
	Rds       *redis.Client
	ErrTo     *util.ErrToJson
	Db        *sql.DB
}

func (d Developers) GetPath() string {
	return "developers"
}

func (d Developers) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}
	var err error
	var result *sql.Rows
	var marshaled []byte

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	result, err = d.Db.Query(repository.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	dev := entity.Developer{}
	result.Next()
	err = result.Scan(&dev.Name, &dev.Country, &dev.SteamId)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	marshaled, err = json.Marshal(dev)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (d Developers) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	result, err = d.Db.Query(repository.SelectDevelopers)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	devs := make([]entity.Developer, 0)
	dev := entity.Developer{}
	for result.Next() {
		err = result.Scan(&dev.Name, &dev.Country, &dev.SteamId)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		devs = append(devs, dev)
	}
	marshaled, err = json.Marshal(devs)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (d Developers) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	var response *SteamResponse
	response = &SteamResponse{GameList: make(map[string]SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + idS)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	devNameArr := response.GameList[idS].Data.Developers
	devCountry := randomdata.Country(randomdata.FullCountry)
	if len(devNameArr) > 0 {
		devName := strings.Join(devNameArr, " ")
		_, err = d.Db.Query("INSERT INTO developers (name, country, steam_id) VALUES ($1, $2, $3)", devName, devCountry, id)
		if err != nil {
			if strings.Contains(err.Error(), "\"devsteam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				d.Logger.Error(errors.New("409 - Developer already exists"))
				d.ErrTo.ErrToJson(w, errors.New("409 - Developer already exists"))
				return
			} else {
				d.Logger.Error(err)
				d.ErrTo.ErrToJson(w, util.ErrSWW)
				return
			}
		}
		d.Rds.Del(context.Background(), DevelopersCounter)
		d.Logger.Infof("Dev with id: %s added. %s is flushed", id, DevelopersCounter)
		_, err = w.Write([]byte(`{"msg":"Success"}`))
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		d.Logger.Error(errors.New("409 - No developer with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - No developer with such id"))
	}
}
func (d Developers) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var err error
	var response *sql.Rows

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	response, err = d.Db.Query(repository.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	if !response.Next() {
		d.Logger.Error(errors.New("409 - no developer to delete with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - no developer to delete with such id"))
		return
	}
	_, err = d.Db.Query(repository.DeleteDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	d.Rds.Del(context.Background(), DevelopersCounter)
	d.Logger.Infof("Dev with id: %s deleted. %s is flushed", id, DevelopersCounter)
	_, err = w.Write([]byte(`{"msg":"Success"}`))
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (d Developers) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	devStruct := DevPubIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	err = json.Unmarshal(requestBody, &devStruct)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	if err = d.Validator.Struct(&devStruct); err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	rStruct := reflect.ValueOf(devStruct)
	fieldTag := reflect.TypeOf(devStruct)
	var fields []interface{}
	k := 2
	fields = append(fields, idS)

	for i := 0; i < rStruct.NumField(); i++ {
		if !rStruct.Field(i).IsNil() {
			repository.UpdateDeveloperById += fmt.Sprintf("%s = $%d, ", fieldTag.Field(i).Tag.Get("json"), k)
			fields = append(fields, rStruct.Field(i).Interface())
			k++
		}
	}
	repository.UpdateDeveloperById = repository.UpdateDeveloperById[:len(repository.UpdateDeveloperById)-2]
	repository.UpdateDeveloperById += " WHERE steam_id = $1"

	response, err = d.Db.Query(repository.SelectDeveloperById, id)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	if !response.Next() {
		d.Logger.Error(errors.New("409 - no developer to update with such id"))
		d.ErrTo.ErrToJson(w, errors.New("409 - no developer to update with such id"))
		return
	}
	_, err = d.Db.Exec(repository.UpdateDeveloperById, fields...)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write([]byte(`{"msg":"Success"}`))
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}

func (d Developers) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		d.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}
	count := Counter{}
	var marshaled []byte
	var err error
	count.Count, err = d.Rds.Get(context.Background(), DevelopersCounter).Result()
	if errors.Is(err, redis.Nil) {
		err = d.Db.QueryRow(repository.GetDevelopersCount).Scan(&count.Count)
		d.Logger.Infof("Redis key %s is empty, getting data from DB", DevelopersCounter)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		d.Rds.Append(context.Background(), DevelopersCounter, count.Count)
		marshaled, err = json.Marshal(count)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			d.Logger.Error(err)
			d.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		return
	}

	d.Logger.Infof("Redis key %s is found, getting data from Redis", DevelopersCounter)
	marshaled, err = json.Marshal(count)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		d.Logger.Error(err)
		d.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	return
}
