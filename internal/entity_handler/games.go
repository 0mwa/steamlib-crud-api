package entity_handler

import (
	"TestProject/internal"
	"TestProject/internal/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Games struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
}

func (g Games) errToJson(w http.ResponseWriter, externalError error) {
	errrr := internal.ErrOut{externalError.Error()}
	marshaled, err := json.Marshal(errrr)
	if err != nil {
		g.Logger.Error(err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		return
	}
}

func (g Games) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var result *sql.Rows
	var marshaled []byte
	id := r.PathValue("id")
	db := internal.GetBD()
	result, err = db.Query(internal.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	game := entity.Game{}
	result.Next()
	err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	marshaled, err = json.Marshal(game)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
}
func (g Games) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.errToJson(w, errors.New(MethodError))
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	db := internal.GetBD()
	result, err = db.Query(internal.SelectGames)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}

	games := make([]entity.Game, 0)
	game := entity.Game{}
	for result.Next() {
		err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
		if err != nil {
			g.Logger.Error(err)
			g.errToJson(w, err)
			return
		}
		games = append(games, game)
	}
	marshaled, err = json.Marshal(games)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
}
func (g Games) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		g.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.errToJson(w, errors.New(MethodError))
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	var response *internal.SteamResponse
	response = &internal.SteamResponse{GameList: make(map[string]internal.SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	gameName := response.GameList[id].Data.Name
	gameImage := response.GameList[id].Data.HeaderImage
	gameDescription := response.GameList[id].Data.ShortDescription[:internal.MyMin(255, cap([]byte(response.GameList[id].Data.ShortDescription)))]
	if gameName != "" {
		_, err = db.Query("INSERT INTO games (steam_id, name, img, description) VALUES ($1, $2, $3, $4)", id, gameName, gameImage, gameDescription)
		if err != nil {
			if strings.Contains(err.Error(), "\"steam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				g.errToJson(w, errors.New("409 - Game already exists"))
			} else {
				g.Logger.Error(err)
				g.errToJson(w, err)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		g.errToJson(w, errors.New("409 - No game with such id"))
	}
}
func (g Games) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		g.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var response *sql.Rows

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	if !response.Next() {
		g.Logger.Error(err)
		g.errToJson(w, errors.New("409 - no game to delete with such id"))
		return
	}
	_, err = db.Query(internal.DeleteGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
}
func (g Games) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		g.Logger.Error(MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.errToJson(w, errors.New(MethodError))
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	gameStruct := internal.GameIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	err = json.Unmarshal(requestBody, &gameStruct)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}

	if err = g.Validator.Struct(&gameStruct); err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}

	db := internal.GetBD()
	id := r.PathValue("id")
	response, err = db.Query(internal.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}
	if !response.Next() {
		g.Logger.Error(err)
		g.errToJson(w, errors.New("409 - no game to update with such id"))
		return
	}
	_, err = db.Query(internal.UpdateGameById, gameStruct.Name, gameStruct.Img, gameStruct.Description, gameStruct.Rating, gameStruct.DeveloperId, gameStruct.PublisherId, id)
	if err != nil {
		g.Logger.Error(err)
		g.errToJson(w, err)
		return
	}

}
