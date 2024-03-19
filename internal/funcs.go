package internal

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

var db *sql.DB

func GetBD() *sql.DB {
	file := "postgres://postgres:omw@localhost:5432/postgres?sslmode=disable"
	var err error
	if db == nil {
		db, err = sql.Open("postgres", file)
		if err != nil {
			panic(err)
		}
	}
	return db
}

func MyMin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// SteamAPI json response parse

func (r SteamResponse) UnmarshalJSON(data []byte) error {
	elements := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &elements)
	if err != nil {
		panic(err)
	}
	for k, v := range elements {
		element := SteamResponseElement{}
		err = json.Unmarshal(v, &element)
		if err != nil {
			panic(err)
		}
		r.GameList[k] = element
	}
	return nil
}

// SteamAPI json response struct

type SteamResponseElementData struct {
	Name             string `json:"name"`
	HeaderImage      string `json:"header_image"`
	ShortDescription string `json:"short_description"`
}

type SteamResponseElement struct {
	Data SteamResponseElementData `json:"data"`
}

type SteamResponse struct {
	GameList map[string]SteamResponseElement
}

// OurAPI json response struct

type IdIn struct {
	Id string `json:"id"`
}

type GameIn struct {
	Name        *string `json:"name"`
	Img         *string `json:"img"`
	Description *string `json:"description"`
	Rating      *string `json:"rating"`
	DeveloperId *string `json:"developer_id"`
	PublisherId *string `json:"publisher_id"`
	SteamId     *string `json:"steam_id"`
}

type DevIn struct {
	Name    *string `json:"name"`
	Country *string `json:"country"`
}

type GameOut struct {
	Name        *string `json:"name"`
	Img         *string `json:"img"`
	Description *string `json:"description"`
	Rating      *string `json:"rating"`
	DeveloperId *string `json:"developer_id"`
	PublisherId *string `json:"publisher_id"`
	SteamId     *string `json:"steam_id"`
}

type DevOut struct {
	Name    *string `json:"name"`
	Country *string `json:"country"`
}

type ErrOut struct {
	Error string `json:"error"`
}
