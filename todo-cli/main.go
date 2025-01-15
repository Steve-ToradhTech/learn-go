package main

import (
	"bufio"
	"container/list"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const filePath = "db.csv"

type record struct {
	id       string
	item     string
	complete bool
}

func listDisplay(list *list.List) {
	for item := list.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
	}
}

func searchById(list *list.List) {
	for item := list.Front(); item != nil; item = item.Next() {
		list.Remove(item)
	}
}

func loadDBIntoMemory() *list.List {
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
	list := list.New()
	for _, rec := range rawRecords {
		boolConv, err := strconv.ParseBool(rec[1])
		if err != nil {
			log.Fatal(err)
		}
		tempRec := record{item: rec[0], complete: boolConv}
		list.PushBack(tempRec)
	}
	return list
}

func main() {
	// Default output
	fmt.Println("Welcome to todos application!\nUse h for all possible commands.")
	// Initialize the linkedlist from .csv file
	list := loadDBIntoMemory()

	// Read input from user and validate it against a list of possible uses.
	for {
		fmt.Print("flag: ")
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

		switch switchChar {
		case "v":
			listDisplay(list)
		// Take user input and add new task to linked list.
		case "a":
			fmt.Println("Enter task:")
			taskInput := bufio.NewReader(os.Stdin)
			task, err := taskInput.ReadString('\n')
			if err != nil {
				log.Fatal("error reading input:", err)
			}
			id := uuid.New()
			newTask := record{id: id.String(), item: task, complete: false}
			list.PushBack(newTask)
		// Delete an task
		case "d":
			fmt.Print(": ")
			//switchInput := bufio.NewReader(os.Stdin)
			//input, err := switchInput.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			searchById(list)
		case "q":
			fmt.Println("Quit")
		case "default", "h":
			fmt.Println("v: view\na: add\nd: delete\nq: quit")
		}
	}
}
