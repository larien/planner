package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/larien/planner/cli"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cli := cli.New()
	option := os.Args[1]
	switch option {
	case "boards":
		cli.Boards()
	case "lists":
		cli.Lists(cli.Config.BoardID)
	case "week":
		cli.CreateWeek(cli.Config.BoardID)
	default:
		fmt.Fprint(os.Stderr, "unavailable option\n")
	}
}

var usage = `Welcome to Planner CLI!

Usage: plan [options...]

Options:
    boards - Lists all open boards
    lists <board_id> - Lists all cards in a board
	week - Creates board for the entire week
`
