package main

import (
	"fmt"
	"os"

	"github.com/larien/planner/cli"
)

func main() {

	fmt.Println(usage)

	cli := cli.New()
	option := os.Args[1]
	switch option {
	case "boards":
		cli.Boards()
	case "lists":
		id := os.Args[2]
		cli.Lists(id)
	default:
		fmt.Fprint(os.Stderr, "unavailable option\n")
	}
}

var usage = `Welcome to Planner CLI!

Usage: plan [options...]

Options:
    boards - Lists all open boards
    lists <board_id> - Lists all cards in a board
`
