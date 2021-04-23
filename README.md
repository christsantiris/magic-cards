# Set Up

## To set up the DB in Elephant SQL (Postgres): 

`create table cards (id serial, name varchar, color varchar, standard_legal boolean, type varchar, rarity varchar, set varchar, casting_cost int)`
`insert into cards (name, color, standard_legal, type, rarity, set, casting_cost) values ('Goldspan Dragon', 'Red', true, 'Creature', 'Mythic Rare', 'Kaldheim', 5)`
`SELECT * FROM cards` 
