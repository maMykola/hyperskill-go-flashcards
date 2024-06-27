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

func (fc *FlashCards) addCard(term, definition *string) {
	fc.Terms[*term] = definition
	fc.Definitions[*definition] = term
}

func addCard() {
	term := getCardInfo("card", getString("The card:"), &flashcards.Terms)
	definition := getCardInfo("definition", getString("The definition of the card:"), &flashcards.Definitions)

	flashcards.addCard(&term, &definition)

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
	term := getString("Which card?")

	if definition, ok := flashcards.Terms[term]; ok {
		delete(flashcards.Terms, term)
		delete(flashcards.Definitions, *definition)
		fmt.Println("The card has been removed.")
	} else {
		fmt.Printf("Can't remove \"%s\": there is no such card.\n", term)
	}
}

func readLine(scanner *bufio.Scanner) (string, bool) {
	var text string
	var ok = true

	for text == "" && ok {
		if ok = scanner.Scan(); ok {
			text = strings.TrimSpace(scanner.Text())
		}
	}

	return text, ok
}

func importCards() {
	filename := getString("File name:")

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("File not found.")
		return
	}
	defer file.Close()

	numCards := 0
	scanner := bufio.NewScanner(file)
	for {
		term, ok1 := readLine(scanner)
		definition, ok2 := readLine(scanner)

		if !ok1 || !ok2 {
			fmt.Printf("%d cards have been loaded.\n", numCards)
			return
		}

		flashcards.addCard(&term, &definition)
		numCards++
	}
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
