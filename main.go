package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type FlashCards struct {
	Terms       map[string]*string
	Definitions map[string]*string
	Errors      map[string]int
}

var flashcards = FlashCards{
	Terms:       make(map[string]*string),
	Definitions: make(map[string]*string),
	Errors:      make(map[string]int),
}

var buffer strings.Builder

func main() {
	for {
		action := getString("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")

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
		case "log":
			saveLog()
		case "hardest card":
			hardestCard()
		case "reset stats":
			resetStats()
		case "exit":
			display("Bye bye!\n")
			return
		}

		display("\n")
	}
}

func (fc *FlashCards) AddCard(term, definition string, numErrors int) {
	fc.Terms[term] = &definition
	fc.Definitions[definition] = &term
	fc.Errors[term] = numErrors
}

func (fc *FlashCards) RemoveCard(term string) bool {
	if definition, ok := fc.Terms[term]; ok {
		delete(fc.Terms, term)
		delete(fc.Definitions, *definition)
		delete(fc.Errors, *definition)

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
		display("Correct!\n")
		return
	}

	fc.Errors[term]++

	if t, ok := fc.Definitions[definition]; ok {
		display(fmt.Sprintf(
			"Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n",
			*fc.Terms[term],
			*t,
		))
	} else {
		display(fmt.Sprintf("Wrong. The right answer is \"%s\".\n", *fc.Terms[term]))
	}
}

func (fc *FlashCards) ResetStats() {
	for term := range fc.Errors {
		fc.Errors[term] = 0
	}
}

func (fc *FlashCards) GetHardestCards() ([]string, int) {
	var cards = make([]string, 0, len(fc.Errors))
	var maxErrors int

	for term, numErrors := range fc.Errors {
		if numErrors > maxErrors {
			maxErrors = numErrors
			cards = cards[:1]
			cards[0] = term
		} else if numErrors > 0 && numErrors == maxErrors {
			cards = append(cards, term)
		}
	}

	return cards, maxErrors
}

func addCard() {
	term := getCardInfo("card", getString("The card:"), &flashcards.Terms)
	definition := getCardInfo("definition", getString("The definition of the card:"), &flashcards.Definitions)

	flashcards.AddCard(term, definition, 0)

	display(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.\n", term, definition))
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
		display("The card has been removed.\n")
	} else {
		display(fmt.Sprintf("Can't remove \"%s\": there is no such card.\n", term))
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
		display("File not found.\n")
		return
	}
	defer file.Close()

	numCards := 0
	scanner := bufio.NewScanner(file)
	for {
		term, ok1 := readLine(scanner)
		definition, ok2 := readLine(scanner)
		number, ok3 := readLine(scanner)

		if !ok1 || !ok2 || !ok3 {
			display(fmt.Sprintf("%d cards have been loaded.\n", numCards))
			return
		}

		numErrors, _ := strconv.Atoi(number)
		flashcards.AddCard(term, definition, numErrors)
		numCards++
	}
}

func exportCards() {
	filename := getString("File name:")

	file, err := os.Create(filename)
	if err != nil {
		display(errors.Unwrap(err).Error(), "\n")
		return
	}
	defer file.Close()

	for term, definition := range flashcards.Terms {
		file.WriteString(term)
		file.WriteString("\n")
		file.WriteString(*definition)
		file.WriteString("\n")
		file.WriteString(strconv.Itoa(flashcards.Errors[term]))
		file.WriteString("\n\n")
	}

	display(fmt.Sprintf("%d cards have been saved.\n", len(flashcards.Terms)))
}

func quiz() {
	numCards := getInt("How many times to ask?")
	terms := flashcards.RandomTerms(numCards)

	for _, term := range terms {
		flashcards.Check(term)
	}
}

func saveLog() {
	filename := getString("File name:")

	file, err := os.Create(filename)
	if err != nil {
		display(errors.Unwrap(err).Error(), "\n")
		return
	}
	defer file.Close()

	file.WriteString(buffer.String())

	display("The log has been saved.\n")
}

func hardestCard() {
	hardestCards, numErrors := flashcards.GetHardestCards()
	if len(hardestCards) == 0 {
		display("There are no cards with errors.\n")
	} else if len(hardestCards) == 1 {
		display(fmt.Sprintf("The hardest card is \"%s\". You have %d errors answering it.\n", hardestCards[0], numErrors))
	} else {
		display(fmt.Sprintf("The hardest cards are \"%s\".", strings.Join(hardestCards, "\", \"")))
	}
}

func resetStats() {
	flashcards.ResetStats()
	fmt.Println("Card statistics have been reset.")
}

func getString(prompt string) string {
	display(prompt, "\n")
	return strings.TrimSpace(readString())
}

func getInt(prompt string) int {
	input := getString(prompt)
	num, _ := strconv.Atoi(input)
	return num
}

func display(text ...string) {
	for _, t := range text {
		fmt.Print(t)
		buffer.WriteString(t)
	}
}

func readString() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	buffer.WriteString(text)

	return text
}
