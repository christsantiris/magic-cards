# Set Up

## To set up the DB in Elephant SQL (Postgres): 

### `create table cards (id serial, name varchar, color varchar, standard_legal boolean, type varchar, rarity varchar, set varchar, casting_cost int)`
### `insert into cards (name, color, standard_legal, type, rarity, set, casting_cost) values ('Goldspan Dragon', 'Red', true, 'Creature', 'Mythic Rare', 'Kaldheim', 5)`
### `SELECT * FROM cards` 

### To create users table
### `create table users (id serial primary key, email text not null unique, password text not null);`

## To run the application in debug mode: 
### `go get github.com/derekparker/delve/cmd/dlv`
### then hit the play button in vs code

## To create a go module and add dependencies:
### `go mod init <module name> i.e github.com/christsantiris/magic-cards`
### `go mod tidy`

## To build and run the project
### `go build && ./magic-cards`

