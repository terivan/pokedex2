package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	po "pokedex2/internal/PokeAPImanager"
	ca "pokedex2/internal/pokecache"
	"strings"
)


type config struct {
	NextUrl     string
	PreviousUrl string
	cache ca.Cache
}

type cliCommand struct {
	name        string
	description string
	command     func(*config) error
}

// func (c *config) exploreCommandFunc() error {
// 	return nil
// }

func (c *config) mapCommandFunc(forward bool) error {

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

	val, ok := c.cache.Get(urlToUse)
	var res []byte
	var err error
	if ok {
		res = val
		err = nil
	}  else {
		res, err = po.GetLocations(urlToUse)
		if err != nil {
			fmt.Println("Couldn't read locations from the interwebs!")
			return err
		}
	}
	
	var locations po.Locations
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
				help: Gives instructions/help
				exit: Exit the Pokedex
				map: Next 20 cities
				mapb: Previous 20 cities
				explore: Explore pokemon in location`)
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
		description: "Next 20 cities",
		command: func(c *config) error {
			return c.mapCommandFunc(true)
		}}

	commandMap["mapb"] = cliCommand{
		name:        "mapb",
		description: "Previous 20 cities",
		command: func(c *config) error {
			return c.mapCommandFunc(false)
		}}

	// commandMap["explore"] = cliCommand{
	// 	name:        "explore",
	// 	description: "Explore pokemon in location",
	// 	command: func(c *config) error {
	// 		return c.exploreCommandFunc(false)
	// 	}}

	return commandMap
}

func cleanInput(text string) []string {
	lowercase_input := strings.ToLower(text)
	whitespace_removed := strings.Trim(lowercase_input, " ")
	split_string := strings.Fields(whitespace_removed)
	return split_string
}

func main() {

	cfg := config{
		NextUrl: "https://pokeapi.co/api/v2/location-area/?limit=20&offset=20",
		PreviousUrl: "",
		cache: *ca.NewCache(5),
	}



	mapOfFuncs := cfg.commandMap()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

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
