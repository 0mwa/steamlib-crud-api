package entity_handler

import (
	"TestProject/internal"
	"TestProject/internal/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Games struct {
	Logger *zap.SugaredLogger
}

func (g Games) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodGet {
		g.Logger.Error("405 - Wrong method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err = w.Write([]byte("405 - Wrong method\n"))
		if err != nil {
			g.Logger.Error(err)
			return
		}
		return
	}

	var errrr internal.ErrOut
	var result *sql.Rows
	var marshaled []byte

	id := r.PathValue("id")
	//fmt.Printf("\n Got id : %s \n", id)
	db := internal.GetBD()

	result, err = db.Query(internal.SelectGameById, id)
	if err != nil {
		g.Logger.Error(err)
		return
	}
	game := entity.Game{}
	//fmt.Printf("\n %v \n", result)
	result.Next()
	err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.Steam)
	if err != nil {
		errrr = internal.ErrOut{err.Error()}
		marshaled, err = json.Marshal(errrr)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		return
	}
	marshaled, err = json.Marshal(game)
	if err != nil {
		errrr = internal.ErrOut{err.Error()}
		marshaled, err = json.Marshal(errrr)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		return
	}
}
func (g Games) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Wrong method\n"))
		g.Logger.Error("405 - Wrong method")
		return
	}

	var errrr internal.ErrOut
	var marshaled []byte
	var result *sql.Rows
	var getQuery = r.URL.RawQuery
	var err error

	db := internal.GetBD()
	switch getQuery {
	case "sort:name":
		result, err = db.Query(internal.SelectGamesSortName)
		if err != nil {
			g.Logger.Error(err)
			return
		}
	default:
		result, err = db.Query(internal.SelectGames)
		if err != nil {
			g.Logger.Error(err)
			return
		}
	}
	games := make([]entity.Game, 0)
	game := entity.Game{}
	for result.Next() {
		err = result.Scan(&game.Name, &game.Img, &game.Description, &game.Rating, &game.DeveloperId, &game.PublisherId, &game.Steam)
		if err != nil {
			errrr = internal.ErrOut{err.Error()}
			marshaled, err = json.Marshal(errrr)
			if err != nil {
				g.Logger.Error(err)
				return
			}
			_, err = w.Write(marshaled)
			if err != nil {
				g.Logger.Error(err)
				return
			}
			return
		}
		games = append(games, game)
	}
	marshaled, err = json.Marshal(games)
	if err != nil {
		errrr = internal.ErrOut{err.Error()}
		marshaled, err = json.Marshal(errrr)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		_, err = w.Write(marshaled)
		if err != nil {
			g.Logger.Error(err)
			return
		}
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		g.Logger.Error(err)
		return
	}
}
func (g Games) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Wrong method\n"))
		g.Logger.Error("405 - Wrong method")
		return
	}

	db := internal.GetBD()
	//apiKey := os.Getenv("API_KEY")
	id := r.PathValue("id")
	var response *internal.SteamResponse
	response = &internal.SteamResponse{GameList: make(map[string]internal.SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		g.Logger.Error(err)
		return
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		g.Logger.Error(err)
		return
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		g.Logger.Error(err)
		return
	}
	//fmt.Println(*response)
	gameName := response.GameList[id].Data.Name
	gameImage := response.GameList[id].Data.HeaderImage
	gameDescription := response.GameList[id].Data.ShortDescription[:internal.MyMin(255, cap([]byte(response.GameList[id].Data.ShortDescription)))]
	//fmt.Println(gameName)
	if gameName != "" {
		_, err = db.Query("INSERT INTO games (steam_id, name, img, description) VALUES ($1, $2, $3, $4)", id, gameName, gameImage, gameDescription)
		if err != nil {
			if strings.Contains(err.Error(), "\"steam_id_unique\"") {
				fmt.Println(err)
				w.WriteHeader(http.StatusConflict)
				_, err = w.Write([]byte("409 - Game already exists!\n"))
				if err != nil {
					g.Logger.Error(err)
					return
				}
			} else {
				g.Logger.Error(err)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte("409 - No game with such id!\n"))
		if err != nil {
			g.Logger.Error(err)
			return
		}
	}
}
func (g Games) Del(w http.ResponseWriter, r *http.Request) {}
func (g Games) Put(w http.ResponseWriter, r *http.Request) {}
