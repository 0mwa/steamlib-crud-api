package internal

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

var db *sql.DB

const MethodError string = "405 - Wrong method"

// API json err response

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
	Name             string   `json:"name"`
	HeaderImage      string   `json:"header_image"`
	ShortDescription string   `json:"short_description"`
	Publishers       []string `json:"publishers"`
	//Developers     []string `json:"developers"`
}

type SteamResponseElement struct {
	Data SteamResponseElementData `json:"data"`
}

type SteamResponse struct {
	GameList map[string]SteamResponseElement
}

// OurAPI json response struct

type GameIn struct {
	Name        *string `json:"name" validate:"max=255"`
	Img         *string `json:"img" validate:"max=255"`
	Description *string `json:"description" validate:"max=255"`
	Rating      *int    `json:"rating" validate:"gte=0,lte=10"`
	DeveloperId *int    `json:"developer_id" validate:"gte=0"`
	PublisherId *int    `json:"publisher_id" validate:"gte=0"`
}

type DevPubIn struct {
	Name    *string `json:"name" validate:"max=255"`
	Country *string `json:"country" validate:"max=100"`
}

type ErrOut struct {
	Error string `json:"error"`
}

type Counter struct {
	Count string `json:"count"`
}
