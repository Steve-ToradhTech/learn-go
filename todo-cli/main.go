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

const filePath = "db.csv"

type record struct {
	id       int
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

func DeleteNode(head *Node, value int) *Node {
	if head == nil {
		return nil
	}

	// If the head node matches the value, return the next node
	if head.data.id == value {
		return head.next
	}

	// Traverse the list to find the node to delete
	currentNode := head
	for currentNode.next != nil {
		println(value)
		if currentNode.next.data.id == value {
			// Skip the node to delete it
			currentNode.next = currentNode.next.next
			return head
		}
		currentNode = currentNode.next
	}

	// If the value is not found, return the original head
	return head
}

func (list *LinkedList) Display() {
	current := list.head

	if current == nil {
		fmt.Println("Linked list is empty.")
		return
	}

	for current != nil {
		fmt.Printf("%v\n", current.data)
		current = current.next
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
		intConv, err := strconv.ParseInt(rec[0], 0, 16)
		boolConv, err := strconv.ParseBool(rec[2])
		if err != nil {
			log.Fatal(err)
		}
		tempRec := record{id: int(intConv), item: rec[1], complete: boolConv}
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
		loadDBIntoMemory(&dataLinkedList)

		switch switchChar {
		case "v":
			dataLinkedList.Display()
		case "a":
			fmt.Println("Add")
		case "d":
			fmt.Print(": ")
			switchInput := bufio.NewReader(os.Stdin)
			char, _, err := switchInput.ReadRune()
			deleteIDChar := string(char)
			deleteID, err := strconv.ParseInt(deleteIDChar, 0, 16)
			if err != nil {
				log.Fatal(err)
			}
			dataLinkedList.head = DeleteNode(dataLinkedList.head, int(deleteID))
		case "q":
			fmt.Println("Quit")
		case "default", "h":
			fmt.Println("v: view\na: add\nd: delete\nq: quit")
		}
	}
}
