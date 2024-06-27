package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FlashCards struct {
	Terms       map[string]*string
	Definitions map[string]*string
}

var flashcards = FlashCards{
	Terms:       make(map[string]*string),
	Definitions: make(map[string]*string),
}

func main() {
	for {
		action := getString("Input the action (add, remove, import, export, ask, exit):")

		switch action {
		case "add":
			addCard()
		case "remove":
			removeCard()
		case "import":
			importCards()
		case "export":
			exportCards()
		case "ask":
			quiz()
		case "exit":
			fmt.Println("Bye bye!")
			return
		}

		fmt.Println()
	}
}

func addCard() {
	term := getCardInfo("card", getString("The card:"), &flashcards.Terms)
	definition := getCardInfo("definition", getString("The definition of the card:"), &flashcards.Definitions)

	flashcards.Terms[term] = &definition
	flashcards.Definitions[definition] = &term

	fmt.Printf("The pair (\"%s\":\"%s\") has been added.\n", term, definition)
}

func getCardInfo(name, value string, list *map[string]*string) string {
	for {
		if _, ok := (*list)[value]; !ok {
			break
		}
		value = getString(fmt.Sprintf("The %s \"%s\" already exists. Try again:", name, value))
	}
	return value
}

func removeCard() {
	// todo: stub
}

func importCards() {
	// todo: stub
}

func exportCards() {
	// todo: stub
}

func quiz() {
	// todo: stub
}

func getString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(prompt)
	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}
