package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const base = "https://api.trello.com/1"

type CLI struct {
	Config Config
}

func New() CLI {
	var cli CLI

	cli.Config = read()

	return cli
}

func (c *CLI) Boards() {
	type Board struct {
		Name     string `json:"name"`
		Closed   bool   `json:"closed"`
		ID       string `json:"id"`
		ShortURL string `json:"shortUrl"`
	}

	resp, err := http.Get(fmt.Sprintf("%s/members/me/boards?key=%s&token=%s&fields=name,closed,shortUrl", base, c.Config.Key, c.Config.Token))
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

func (c *CLI) Lists(boardID string) {
	type List struct {
		Name    string `json:"name"`
		Closed  bool   `json:"closed"`
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
	}

	resp, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?key=%s&token=%s", base, boardID, c.Config.Key, c.Config.Token))
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

func (c *CLI) CreateWeek(boardID string) {
	fmt.Println("Creating week")
	for i := 7; i > 0; i-- {
		lastDay := time.Now().AddDate(0, 0, i)
		day := lastDay.Day()
		weekday := days[lastDay.Weekday()]
		name := fmt.Sprintf("%d | %s", day, weekday)
		fmt.Println("Creating: ", name)
		if weekday == "Sábado" || weekday == "Domingo" {
			c.CreateWeekend(name, boardID)
			continue
		}
		c.CreateWorkday(name, boardID)
	}
}

func (c *CLI) CreateWeekend(name, boardID string) {
	listID := c.CreateList(name, boardID)
	c.CreateCard(listID, "Exercício físico ~ 30min")
	c.CreateCard(listID, "Café da manhã ~ 30min")
	c.CreateCard(listID, "Almoço ~ 60min")
	c.CreateCard(listID, "Café da tarde ~ 30min")
	c.CreateCard(listID, "Janta ~ 60min")
}

func (c *CLI) CreateWorkday(name, boardID string) {
	listID := c.CreateList(name, boardID)
	c.CreateCard(listID, "Exercício físico ~ 30min")
	c.CreateCard(listID, "Café da manhã ~ 30min")
	c.CreateCard(listID, "Almoço ~ 60min")
	c.CreateCard(listID, "Café da tarde ~ 30min")
	c.CreateCard(listID, "Janta ~ 60min")
	c.CreateCard(listID, "Foco 1 ~ 60min")
	c.CreateCard(listID, "Foco 2 ~ 60min")
	c.CreateCard(listID, "Foco 3 ~ 60min")
	c.CreateCard(listID, "Foco 4 ~ 60min")
	c.CreateCard(listID, "Foco 5 ~ 60min")
	c.CreateCard(listID, "Foco 6 ~ 60min")
}

func (c *CLI) CreateCard(listID, name string) {
	data := url.Values{}
	data.Set("key", c.Config.Key)
	data.Set("token", c.Config.Token)
	data.Set("name", name)
	data.Set("idList", listID)
	_, err := http.Post(fmt.Sprintf("%s/cards?%s", base, data.Encode()), "application/json", nil)
	if err != nil {
		log.Fatalf("failed to post request: %v", err)
	}
}

type Weekday int

const (
	Domingo Weekday = iota
	Segunda
	Terça
	Quarta
	Quinta
	Sexta
	Sábado
)

var days = [...]string{
	"Domingo",
	"Segunda-feira",
	"Terça-feira",
	"Quarta-feira",
	"Quinta-feira",
	"Sexta-feira",
	"Sábado",
}

func (c *CLI) CreateList(name, boardID string) (listID string) {
	type List struct {
		Name    string `json:"name"`
		Closed  bool   `json:"closed"`
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
	}
	data := url.Values{}
	data.Set("key", c.Config.Key)
	data.Set("token", c.Config.Token)
	data.Set("name", name)
	data.Set("idBoard", boardID)
	resp, err := http.Post(fmt.Sprintf("%s/lists?%s", base, data.Encode()), "application/json", nil)
	if err != nil {
		log.Fatalf("failed to post request: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}
	var list List
	if err := json.Unmarshal(body, &list); err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}
	return list.ID
}

type Config struct {
	Key     string
	Token   string
	BoardID string
}

func read() Config {
	key := os.Getenv("TRELLO_KEY")
	if key == "" {
		log.Fatal("Key can't be empty")
	}
	token := os.Getenv("TRELLO_TOKEN")
	if token == "" {
		log.Fatal("Token can't be empty")
	}
	boardID := os.Getenv("TRELLO_BOARD_ID")
	if boardID == "" {
		log.Fatal("Board ID can't be empty")
	}
	return Config{
		Key:     key,
		Token:   token,
		BoardID: boardID,
	}
}
