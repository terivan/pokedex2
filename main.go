package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	po "pokedex2/internal/PokeAPImanager"
	"strings"
	 _ "github.com/lib/pq"
	"github.com/joho/godotenv"
	database "pokedex2/internal/database"
	"database/sql"
)

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	NextUrl     string
	PreviousUrl string
	dbQueries *database.Queries
	// UrlToUse    string
	// StepSize    int64
}

type cliCommand struct {
	name        string
	description string
	command     func(*config) error
}

func (c *config) mapCommandFunc(forward bool) error {

	// step := c.StepSize
	//
	var urlToUse string

	if !forward && c.PreviousUrl == "" {
		fmt.Println("You're on the first page")
		return nil
	} else if !forward && c.PreviousUrl != "" {
		urlToUse = c.PreviousUrl
	}

	if forward && c.NextUrl != "" {
		urlToUse = c.NextUrl
	} else if forward && c.NextUrl == "" {
		fmt.Println("Wrong URL!")
	}

	res, err := po.GetLocations(urlToUse)

	if err != nil {
		fmt.Println("Couldn't read!")
		return err
	}

	var locations Locations
	errLoc := json.Unmarshal(res, &locations)

	if errLoc != nil {
		fmt.Println("Error unmarshaling JSON:", errLoc)
		return errLoc
	}

	for _, j := range locations.Results {
		fmt.Println(j.Name)
	}

	c.NextUrl = locations.Next
	c.PreviousUrl = locations.Previous
	// fmt.Println(c)
	return nil

}

func (c *config) helpCommandFunc() error {
	fmt.Println(`Welcome to the Pokedex!
				Usage:
				help: Displays a help message
				exit: Exit the Pokedex
				map: Next 20 cities
				mapb: Previous 20 cities`)
	return nil
}

func (c *config) exitCommandFunc() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *config) commandMap() map[string]cliCommand {
	commandMap := make(map[string]cliCommand)

	commandMap["help"] = cliCommand{
		name:        "help",
		description: "Gives instructions/help",
		command: func(c *config) error {
			return c.helpCommandFunc()
		}}

	commandMap["exit"] = cliCommand{
		name:        "exit",
		description: "Exits the app",
		command: func(c *config) error {
			return c.exitCommandFunc()
		}}

	commandMap["map"] = cliCommand{
		name:        "map",
		description: "Map",
		command: func(c *config) error {
			return c.mapCommandFunc(true)
		}}

	commandMap["mapb"] = cliCommand{
		name:        "mapb",
		description: "Map back",
		command: func(c *config) error {
			return c.mapCommandFunc(false)
		}}

	return commandMap
}

func cleanInput(text string) []string {
	lowercase_input := strings.ToLower(text)
	whitespace_removed := strings.Trim(lowercase_input, " ")
	split_string := strings.Fields(whitespace_removed)
	return split_string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	var cfg config

	cfg.NextUrl = "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0"
	cfg.PreviousUrl = ""
	cfg.dbQueries = dbQueries

	for {
		fmt.Print("Pokedex > ")
		mapOfFuncs := cfg.commandMap()
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		commandText := cleanInput(strings.Split(text, ">")[0])

		if len(commandText) == 0 {
			fmt.Printf("Empty input!\n")
			continue
		}
		commandClean := commandText[0]
		inputCommand, exists := mapOfFuncs[commandClean]

		if exists {
			inputCommand.command(&cfg)
		} else {
			fmt.Printf("Command %v doesn't exist!\n", text)
			continue
		}

	}
}
