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
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
	"strconv"
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
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var err error
	var result *sql.Rows
	var marshaled []byte
	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	result, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	game := entity.Game{}
	result.Next()
	err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	marshaled, err = json.Marshal(game)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (g Games) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var marshaled []byte
	var result *sql.Rows
	var err error

	result, err = g.Db.Query(repository.SelectGames)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	games := make([]entity.Game, 0)
	game := entity.Game{}
	for result.Next() {
		err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.SteamId)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		games = append(games, game)
	}
	marshaled, err = json.Marshal(games)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (g Games) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	var response *SteamResponse
	response = &SteamResponse{GameList: make(map[string]SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + idS)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	gameName := response.GameList[idS].Data.Name
	gameImage := response.GameList[idS].Data.HeaderImage
	gameDescription := response.GameList[idS].Data.ShortDescription[:min(255, cap([]byte(response.GameList[idS].Data.ShortDescription)))]
	if gameName != "" {
		var stmt *sql.Stmt
		stmt, err = g.Db.Prepare("INSERT INTO games (steam_id, name, img, description) VALUES ($1, $2, $3, $4)")
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		_, err = stmt.Query(id, gameName, gameImage, gameDescription)
		if err != nil {
			if strings.Contains(err.Error(), "\"steam_id_unique\"") {
				w.WriteHeader(http.StatusConflict)
				g.Logger.Error(errors.New("409 - Game already exists"))
				g.ErrTo.ErrToJson(w, errors.New("409 - Game already exists"))
				return
			} else {
				g.Logger.Error(err)
				g.ErrTo.ErrToJson(w, util.ErrSWW)
				return
			}
		}
		g.Rds.Del(context.Background(), GamesCounter)
		g.Logger.Infof("Game with id: %s added. %s is flushed", id, GamesCounter)
		_, err = w.Write([]byte(`{"msg":"Success"}`))
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		g.Logger.Error(errors.New("409 - No game with such id"))
		g.ErrTo.ErrToJson(w, errors.New("409 - No game with such id"))
	}
}
func (g Games) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var err error
	var response *sql.Rows

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	response, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
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
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	g.Rds.Del(context.Background(), GamesCounter)
	g.Logger.Infof("Game with id: %s deleted. %s is flushed", id, GamesCounter)
	_, err = w.Write([]byte(`{"msg":"Success"}`))
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}
func (g Games) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
		return
	}

	var err error
	var requestBody []byte
	var response *sql.Rows

	idS := r.PathValue("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrNI)
			return
		}
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	gameStruct := GameIn{}
	requestBody, err = io.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	err = json.Unmarshal(requestBody, &gameStruct)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	if err = g.Validator.Struct(&gameStruct); err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	rStruct := reflect.ValueOf(gameStruct)
	fieldTag := reflect.TypeOf(gameStruct)
	var fields []interface{}
	k := 2
	fields = append(fields, idS)

	for i := 0; i < rStruct.NumField(); i++ {
		if !rStruct.Field(i).IsNil() {
			repository.UpdateGameById += fmt.Sprintf("%s = $%d, ", fieldTag.Field(i).Tag.Get("json"), k)
			fields = append(fields, rStruct.Field(i).Interface())
			k++
		}
	}
	repository.UpdateGameById = repository.UpdateGameById[:len(repository.UpdateGameById)-2]
	repository.UpdateGameById += " WHERE steam_id = $1"

	response, err = g.Db.Query(repository.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	if !response.Next() {
		g.Logger.Error(errors.New("409 - no game to update with such id"))
		g.ErrTo.ErrToJson(w, errors.New("409 - no game to update with such id"))
		return
	}

	var stmt *sql.Stmt
	stmt, err = g.Db.Prepare(repository.UpdateGameById)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	_, err = stmt.Query(fields...)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}

	repository.UpdateGameById = "UPDATE games SET "

	_, err = w.Write([]byte(`{"msg":"Success"}`))
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
}

func (g Games) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		g.Logger.Error(util.ErrMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
		g.ErrTo.ErrToJson(w, util.ErrMethod)
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
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		g.Rds.Append(context.Background(), GamesCounter, count.Count)
		marshaled, err = json.Marshal(count)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			g.Logger.Error(err)
			g.ErrTo.ErrToJson(w, util.ErrSWW)
			return
		}
		return
	}
	g.Logger.Infof("Redis key %s is found, getting data from Redis", GamesCounter)
	marshaled, err = json.Marshal(count)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		g.ErrTo.ErrToJson(w, util.ErrSWW)
		return
	}
	return
}
