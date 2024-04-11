package entity_handler

import "encoding/json"

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
