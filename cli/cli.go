package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const base = "https://api.trello.com/1/"

type CLI struct {
	Auth Auth
}

func New() CLI {
	var cli CLI

	cli.Auth = read()

	return cli
}

func (c *CLI) Boards() {
	type Board struct {
		Name     string `json:"name"`
		Closed   bool   `json:"closed"`
		ID       string `json:"id"`
		ShortURL string `json:"shortUrl"`
	}

	resp, err := http.Get(fmt.Sprintf("%s/members/me/boards?key=%s&token=%s&fields=name,closed,shortUrl", base, c.Auth.Key, c.Auth.Token))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var boards []Board
	if err := json.Unmarshal(body, &boards); err != nil {
		log.Fatalln(err)
	}

	for _, board := range boards {
		if !board.Closed {
			fmt.Printf("[%s] %s: %s\n", board.ID, board.Name, board.ShortURL)
		}
	}
}

func (c *CLI) Lists(id string) {
	type List struct {
		Name    string `json:"name"`
		Closed  bool   `json:"closed"`
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
	}

	resp, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?key=%s&token=%s", base, id, c.Auth.Key, c.Auth.Token))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var lists []List
	if err := json.Unmarshal(body, &lists); err != nil {
		log.Fatalln(err)
	}

	for _, list := range lists {
		if !list.Closed {
			fmt.Printf("[%s][%s] %s\n", list.IDBoard, list.ID, list.Name)
		}
	}
}

type Auth struct {
	Key   string
	Token string
}

func read() Auth {
	key := os.Getenv("TRELLO_KEY")
	if key == "" {
		log.Fatal("Key can't be empty")
	}
	token := os.Getenv("TRELLO_TOKEN")
	if token == "" {
		log.Fatal("Token can't be empty")
	}
	return Auth{
		Key:   key,
		Token: token,
	}
}
