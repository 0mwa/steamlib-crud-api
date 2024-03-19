package entity_handler

import (
	"TestProject/internal"
	"TestProject/internal/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type Games struct {
}

func (Games) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Wrong method"))

		// ToDo: add err log

		return
	}
	var err error
	var result *sql.Rows

	id := r.PathValue("id")
	//fmt.Printf("\n Got id : %s \n", id)
	db := internal.GetBD()

	result, err = db.Query(internal.SelectGameById, id)
	if err != nil {
		panic(err)
	}
	game := entity.Game{}
	fmt.Printf("\n %v \n", result)
	result.Next()
	err = result.Scan(&game.Name, &game.Img, &game.Rating, &game.Description)
	if err != nil {
		io.WriteString(w, "no such id")
	} else {
		var tmplFile = "templates/game.tmpl"
		tmpl, err := template.New("game.tmpl").ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, game)
		if err != nil {
			panic(err)
		}
	}
}
func (Games) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Wrong method"))

		// ToDo: add err log

		return
	}

	var result *sql.Rows
	var getQuery = r.URL.RawQuery
	var err error
	db := internal.GetBD()
	switch getQuery {
	case "sort:name":
		result, err = db.Query(internal.SelectGamesSortName)
		if err != nil {
			panic(err)
		}
	default:
		result, err = db.Query(internal.SelectGames)
		if err != nil {
			panic(err)
		}
	}
	games := make([]entity.Game, 0)
	game := entity.Game{}
	for result.Next() {
		err = result.Scan(&game.Name, &game.Img, &game.Rating, &game.Description)
		if err != nil {
			panic(err)
		}
		games = append(games, game)
	}
	var tmplFile = "templates/games.tmpl"
	tmpl, err := template.New("games.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, games)
	if err != nil {
		panic(err)
	}
}
func (Games) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Wrong method"))

		// ToDo: add err log

		return
	}

	db := internal.GetBD()
	//apiKey := os.Getenv("API_KEY")
	id := r.PathValue("id")
	var response *internal.SteamResponse
	response = &internal.SteamResponse{GameList: make(map[string]internal.SteamResponseElement)}
	get, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		panic(err)
	}
	readall, err := io.ReadAll(get.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(readall, response)
	if err != nil {
		panic(err)
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
				_, err = w.Write([]byte("409 - Game already exists!"))
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte("409 - No game with such id!"))
	}
}
func (Games) Del(w http.ResponseWriter, r *http.Request) {}
func (Games) Put(w http.ResponseWriter, r *http.Request) {}
