# planner

Planner is a CLI to automate daily and weekly planning in Trello, the app I use the most for daily planning and organization.

## Features

- [X] List all boards
- [X] List all lists in a board
- [X] Create a list in a board
- [ ] List all cards in a list
- [X] Create a card in a list
- [ ] Add tags to a card
- [ ] Add a checklist to a card
- [X] Plan week

## Commands

- Apply .env

``` bash
export $(egrep -v '^#' .env | xargs)
```