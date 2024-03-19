package entity

type Game struct {
	Id          int
	Name        string
	Img         string
	Rating      int
	Description string
	Genres      string
	Developer   Developer
	Publisher   Publisher
}
