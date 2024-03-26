package entity

type Game struct {
	Id          int
	Name        string
	Img         string
	Rating      int
	Description string
	Genres      string
	DeveloperId int
	Developer   Developer
	PublisherId int
	Publisher   Publisher
	SteamId     int
}
