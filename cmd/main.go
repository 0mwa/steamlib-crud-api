package main

import (
	"TestProject/internal/entity_handler"
	_ "golang.org/x/net/html"
	"net/http"
)

func main() {

	var a entity_handler.Games

	http.HandleFunc("/games", a.GetAll)
	http.HandleFunc("/games/{id}", a.Get)
	http.HandleFunc("/games/add/{id}", a.Post)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}

//func addGameById(w http.ResponseWriter, r *http.Request) {
//
//	//var result *sql.Rows
//	var err error
//	id := r.PathValue("id")
//	db := internal.GetBD()
//
//	resp, err := http.Get("https://store.steampowered.com/app/" + id)
//	if err != nil {
//		panic(err)
//	}
//	respParse, err := html.Parse(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	buf, err := getHTMLrespGameName(respParse)
//	if err != nil {
//		panic(err)
//	}
//	respstring := collectText(buf)
//	//fmt.Printf("%+v\n\n", buf)
//
//	_, err = db.Query("INSERT INTO games (steam_id, name) VALUES ($1, $2)", id, respstring)
//	if err != nil {
//		fmt.Println(err)
//		w.WriteHeader(http.StatusConflict)
//		_, err = w.Write([]byte("409 - Game already exists!"))
//		if err != nil {
//			panic(err)
//		}
//	}
//}
//
//func getHTMLrespGameName(doc *html.Node) (*html.Node, error) {
//	var appHubAppName *html.Node
//	var crawler func(*html.Node)
//	crawler = func(node *html.Node) {
//		for _, attribute := range node.Attr {
//			if node.Type == html.ElementNode && attribute.Key == "id" && attribute.Val == "appHubAppName" {
//				appHubAppName = node
//				return
//			}
//		}
//		for child := node.FirstChild; child != nil; child = child.NextSibling {
//			crawler(child)
//		}
//	}
//	crawler(doc)
//	if appHubAppName != nil {
//		return appHubAppName, nil
//	}
//	return nil, errors.New("Missing <appHubAppName> in the node tree")
//}

//func collectText(n *html.Node) string {
//	if n.Type == html.TextNode {
//		return n.Data
//	}
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		text := collectText(c)
//		if text != "" {
//			return text
//		}
//	}
//	return ""
//}
