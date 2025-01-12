package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type record struct {
	item     string
	complete bool
}

type Node struct {
	data record
	next *Node
}

type LinkedList struct {
	head *Node
}

func (list *LinkedList) Insert(data record) {
	newNode := &Node{data: data}

	if list.head == nil {
		list.head = newNode
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}

func (list *LinkedList) Display() {
	current := list.head

	if current == nil {
		fmt.Println("Linked list is empty.")
		return
	}

	for current != nil {
		fmt.Printf("%v", current.data)
		current = current.next
	}
}

const filePath = "db.csv"

func readSwitch(input string, dataLinkedList LinkedList) {
	switch input {
	case "v":
		fmt.Println("Linkedlist logic.")
		loadDBIntoMemory(&dataLinkedList)
		dataLinkedList.Display()
	case "a":
		fmt.Println("Add")
	case "d":
		fmt.Println("Delete")
	case "q":
		fmt.Println("Quit")
	case "default", "h":
		fmt.Println("v: view\na: add\nd: delete\nq: quit")
	}
}

func viewDB() {
	// open file throw error if unable to open or read
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer file.Close()

	// read all items from csv file
	csvReader := csv.NewReader(file)
	rawRecords, err := csvReader.ReadAll()
	records := []record{}
	if err != nil {
		log.Fatal("unable to parse file as CSV for "+filePath, err)
	}
	for _, line := range rawRecords {
		boolConv, err := strconv.ParseBool(line[1])
		if err != nil {
			log.Fatal(err)
		}
		tempRec := record{item: line[0], complete: boolConv}
		records = append(records, tempRec)
	}
	printRecords(records)
}

func printRecords(records []record) {
	for _, rec := range records {
		fmt.Printf("%s - %t\n", rec.item, rec.complete)
	}
}

func loadDBIntoMemory(dataLinkedList *LinkedList) {
	// Open the file and defer it's closing until the function exits
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read all data from the csv file and create a linked list
	csvReader := *csv.NewReader(file)
	rawRecords, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, rec := range rawRecords {
		boolConv, err := strconv.ParseBool(rec[1])
		if err != nil {
			log.Fatal(err)
		}
		tempRec := record{item: rec[0], complete: boolConv}
		dataLinkedList.Insert(tempRec)
	}
}

func main() {
	// Default output
	fmt.Println("Welcome to todos application!\nUse h for all possible commands.")

	// Read input from user and validate it against a list of possible uses.
	for {
		fmt.Print(": ")
		switchInput := bufio.NewReader(os.Stdin)
		char, _, err := switchInput.ReadRune()
		switchChar := string(char)
		switchChar = strings.ToLower(switchChar)
		if err != nil {
			log.Fatal(err)
		}
		if switchChar == "q" {
			break
		}

		// Initialize the linked list from file
		dataLinkedList := *&LinkedList{}

		// Logic based on users input
		readSwitch(switchChar, dataLinkedList)
	}
}
