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
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Games struct {
	Logger    *zap.SugaredLogger
	Validator *validator.Validate
	Rds       *redis.Client
	ErrTo     *util.ErrToJson
	Db        *sql.DB
}

func (g Games) GetPath() string {
	return "games"
}

func (g Games) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var err error
	var result *sql.Rows
	var marshaled []byte
	id := r.PathValue("id")
	result, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	game := entity.Game{}
	result.Next()
	err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	marshaled, err = json.Marshal(game)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
}
func (g Games) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	result, err = g.Db.Query(repository.SelectGames)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}

	games := make([]entity.Game, 0)
	game := entity.Game{}
	for result.Next() {
		err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, err)
			return
		}
		games = append(games, game)
	}
	marshaled, err = json.Marshal(games)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
}
func (g Games) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	id := r.PathValue("id")
	var response *SteamResponse
	response = &SteamResponse{GameList: make(map[string]SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	gameName := response.GameList[id].Data.Name
	gameImage := response.GameList[id].Data.HeaderImage
	gameDescription := response.GameList[id].Data.ShortDescription[:util.IntMin(255, cap([]byte(response.GameList[id].Data.ShortDescription)))]
	if gameName != "" {
		_, err = g.Db.Query("INSERT INTO games (steam_id, name, img, description) VALUES ($1, $2, $3, $4)", id, gameName, gameImage, gameDescription)
		if err != nil {
			if strings.Contains(err.Error(), "\"steam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				g.Logger.Error(errors.New("409 - Game already exists"))
				g.ErrTo.ErrToJson(w, errors.New("409 - Game already exists"))
			} else {
				g.Logger.Error(err)
				g.ErrTo.ErrToJson(w, err)
				return
			}
		}
		g.Rds.Del(context.Background(), GamesCounter)
		g.Logger.Infof("Game with id: %s added. %s is flushed", id, GamesCounter)
	} else {
		w.WriteHeader(http.StatusConflict)
		g.Logger.Error(errors.New("409 - No game with such id"))
		g.ErrTo.ErrToJson(w, errors.New("409 - No game with such id"))
	}
}
func (g Games) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var err error
	var response *sql.Rows

	id := r.PathValue("id")
	response, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		g.Logger.Error(errors.New("409 - no game to delete with such id"))
		g.ErrTo.ErrToJson(w, errors.New("409 - no game to delete with such id"))
		return
	}
	_, err = g.Db.Query(repository.DeleteGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	g.Rds.Del(context.Background(), GamesCounter)
	g.Logger.Infof("Game with id: %s deleted. %s is flushed", id, GamesCounter)
}
func (g Games) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	gameStruct := GameIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	err = json.Unmarshal(requestBody, &gameStruct)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}

	if err = g.Validator.Struct(&gameStruct); err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}

	id := r.PathValue("id")
	response, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	if !response.Next() {
		g.Logger.Error(errors.New("409 - no game to update with such id"))
		g.ErrTo.ErrToJson(w, errors.New("409 - no game to update with such id"))
		return
	}
	_, err = g.Db.Query(repository.UpdateGameById, gameStruct.Name, gameStruct.Img, gameStruct.Description, gameStruct.Rating, gameStruct.DeveloperId, gameStruct.PublisherId, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
}

func (g Games) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(handler.MethodError)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, errors.New(handler.MethodError))
		return
	}
	count := Counter{}
	var marshaled []byte
	var err error
	count.Count, err = g.Rds.Get(context.Background(), GamesCounter).Result()
	if errors.Is(err, redis.Nil) {
		err = g.Db.QueryRow(repository.GetGamesCount).Scan(&count.Count)
		g.Logger.Infof("Redis key %s is empty, getting data from DB", GamesCounter)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, err)
			return
		}
		g.Rds.Append(context.Background(), GamesCounter, count.Count)
		marshaled, err = json.Marshal(count)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, err)
			return
		}
		return
	}
	g.Logger.Infof("Redis key %s is found, getting data from Redis", GamesCounter)
	marshaled, err = json.Marshal(count)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, err)
		return
	}
	return
}
