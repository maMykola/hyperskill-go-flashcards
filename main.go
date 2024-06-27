package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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

func (fc *FlashCards) AddCard(term, definition string) {
	fc.Terms[term] = &definition
	fc.Definitions[definition] = &term
}

func (fc *FlashCards) RemoveCard(term string) bool {
	if definition, ok := fc.Terms[term]; ok {
		delete(fc.Terms, term)
		delete(fc.Definitions, *definition)

		return true
	}

	return false
}

func (fc *FlashCards) AllTerms() []string {
	all := make([]string, 0, len(fc.Terms))
	for term := range fc.Terms {
		all = append(all, term)
	}
	return all
}

func (fc *FlashCards) RandomTerms(num int) []string {
	randomTerms := make([]string, 0, num)
	allTerms := fc.AllTerms()
	total := len(allTerms)

	for i := 0; i < num; i++ {
		randomTerms = append(randomTerms, allTerms[rand.Intn(total)])
	}

	return randomTerms
}

func (fc *FlashCards) Check(term string) {
	definition := getString(fmt.Sprintf("Print the definition of \"%s\":", term))

	if *fc.Terms[term] == definition {
		fmt.Println("Correct!")
	} else if t, ok := fc.Definitions[definition]; ok {
		fmt.Printf(
			"Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n",
			*fc.Terms[term],
			*t,
		)
	} else {
		fmt.Printf("Wrong. The right answer is \"%s\".\n", *fc.Terms[term])
	}
}

func addCard() {
	term := getCardInfo("card", getString("The card:"), &flashcards.Terms)
	definition := getCardInfo("definition", getString("The definition of the card:"), &flashcards.Definitions)

	flashcards.AddCard(term, definition)

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

	if flashcards.RemoveCard(term) {
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

		flashcards.AddCard(term, definition)
		numCards++
	}
}

func exportCards() {
	filename := getString("File name:")

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(errors.Unwrap(err))
		return
	}
	defer file.Close()

	for term, definition := range flashcards.Terms {
		file.WriteString(term)
		file.WriteString("\n")
		file.WriteString(*definition)
		file.WriteString("\n\n")
	}

	fmt.Printf("%d cards have been saved.\n", len(flashcards.Terms))
}

func quiz() {
	numCards := getInt("How many times to ask?")
	terms := flashcards.RandomTerms(numCards)

	for _, term := range terms {
		flashcards.Check(term)
	}
}

func getString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(prompt)
	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

func getInt(prompt string) int {
	var num int

	fmt.Println(prompt)
	fmt.Scanln(&num)

	return num
}
