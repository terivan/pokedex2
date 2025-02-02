package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	command     func()
}

func helpCommandFunc() {
	fmt.Println(`Welcome to the Pokedex!
Usage:
help: Displays a help message
exit: Exit the Pokedex`)
}

func exitCommandFunc() {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandMap() map[string]cliCommand {
	commandMap := make(map[string]cliCommand)

	commandMap["help"] = cliCommand{
		name:        "help",
		description: "Gives instructions/help",
		command:     helpCommandFunc}

	commandMap["exit"] = cliCommand{
		name:        "exit",
		description: "Exits the app",
		command:     exitCommandFunc}

	// commandMap["map"] = cliCommand{
	// 	name: "map",
	// 	command: cfg.mapCommandFunc}

	return commandMap
}

func cleanInput(text string) []string {
	lowercase_input := strings.ToLower(text)
	whitespace_removed := strings.Trim(lowercase_input, " ")
	split_string := strings.Fields(whitespace_removed)
	return split_string
}

func main() {
	for {
		fmt.Print("Pokedex > ")
		mapOfFuncs := commandMap()
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
			inputCommand.command()
		} else {
			fmt.Printf("Command %v doesn't exist!\n", text)
			continue
		}

	}
}
