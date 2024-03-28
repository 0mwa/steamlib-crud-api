package entity_handler

import (
	"TestProject/internal"
	"TestProject/internal/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Publishers struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
}

func (p Publishers) GetPath() string {
	return "publishers"
}

func (p Publishers) errToJson(w http.ResponseWriter, externalError error) {
	errrr := internal.ErrOut{externalError.Error()}
	marshaled, err := json.Marshal(errrr)
	if err != nil {
		p.Logger.Error(err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		return
	}
}

func (p Publishers) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.errToJson(w, errors.New(MethodError))
		return
	}
	var err error
	var result *sql.Rows
	var marshaled []byte

	id := r.PathValue("id")
	db := internal.GetBD()
	result, err = db.Query(internal.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	pub := entity.Publisher{}
	result.Next()
	err = result.Scan(&pub.Name, &pub.Country, &pub.SteamId)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	marshaled, err = json.Marshal(pub)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
}
func (p Publishers) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.errToJson(w, errors.New(MethodError))
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	db := internal.GetBD()
	result, err = db.Query(internal.SelectPublishers)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	pubs := make([]entity.Publisher, 0)
	pub := entity.Publisher{}
	for result.Next() {
		err = result.Scan(&pub.Name, &pub.Country, &pub.SteamId)
		if err != nil {
			p.Logger.Error(err)
			p.errToJson(w, err)
			return
		}
		pubs = append(pubs, pub)
	}
	marshaled, err = json.Marshal(pubs)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
}
func (p Publishers) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.errToJson(w, errors.New(MethodError))
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	var response *internal.SteamResponse
	response = &internal.SteamResponse{GameList: make(map[string]internal.SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	pubNameArr := response.GameList[id].Data.Publishers
	pubCountry := randomdata.Country(randomdata.FullCountry)
	if len(pubNameArr) > 0 {
		pubName := strings.Join(pubNameArr, " ")
		_, err = db.Query("INSERT INTO publishers (name, country, steam_id) VALUES ($1, $2, $3)", pubName, pubCountry, id)
		if err != nil {
			if strings.Contains(err.Error(), "\"pubsteam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				p.errToJson(w, errors.New("409 - Publisher already exists"))
			} else {
				p.Logger.Error(err)
				p.errToJson(w, err)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		p.errToJson(w, errors.New("409 - No publisher with such id"))
	}
}
func (p Publishers) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		p.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var response *sql.Rows

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	if !response.Next() {
		p.Logger.Error(err)
		p.errToJson(w, errors.New("409 - no publisher to delete with such id"))
		return
	}
	_, err = db.Query(internal.DeletePublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
}
func (p Publishers) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		p.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		p.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	pubStruct := internal.DevPubIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	err = json.Unmarshal(requestBody, &pubStruct)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}

	if err = p.Validator.Struct(&pubStruct); err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectPublisherById, id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
	if !response.Next() {
		p.Logger.Error(err)
		p.errToJson(w, errors.New("409 - no publisher to update with such id"))
		return
	}
	_, err = db.Query(internal.UpdatePublisherById, pubStruct.Name, pubStruct.Country, id)
	if err != nil {
		p.Logger.Error(err)
		p.errToJson(w, err)
		return
	}
}
