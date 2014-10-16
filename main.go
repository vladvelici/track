package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	STORAGE_PATH = "./data.json"
)

func confirm(msg string) bool {
	fmt.Print(msg)
	fmt.Print(" Continue? y/n: ")
	var result string
	_, err := fmt.Scanln(&result)
	if err != nil {
		return false
	}
	result = strings.ToLower(result)
	switch result {
	case "y", "yes", "sure", "yeah":
		return true
	default:
		return false
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Should have a flag.")
		os.Exit(1)
	}

	// commands that create/delete the index: "rm", "init"
	switch os.Args[1] {
	case "delete":
		if !confirm("This will delete the current index.") {
			return
		}
		if err := deleteIndex(STORAGE_PATH); err != nil {
			deleteIndex(STORAGE_PATH)
		}
		return
	case "init":
		if !confirm("This will destroy all your current data.") {
			return
		}
		if _, err := createIndex(STORAGE_PATH); err != nil {
			fmt.Printf("Cannot create index at %s. %s\n", STORAGE_PATH, err)
		}
		return
	}

	index, err := readIndex(STORAGE_PATH)
	if err != nil {
		fmt.Println("Error reading the index file at %s. %s", STORAGE_PATH, err)
		return
	}

	write := false
	switch os.Args[1] {
	case "work":
		fallthrough
	case "start":
		if len(os.Args) != 3 {
			fmt.Println("Please use `track start <project-name>`")
			return
		}
		write, err = index.StartWorking(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
	case "stop":
		if len(os.Args) != 3 {
			fmt.Println("Please use `track stop <project-name>`")
			return
		}
		write, err = index.StopWorking(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
	case "add":
		write, err = index.AddProjects(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "rm":
		if confirm(fmt.Sprintf("Delete projects %s?", os.Args[2:])) {
			write = rm(index, os.Args[1:])
		}
	case "status":
		status(index, os.Args[2:])
	default:
		fmt.Println("Unkown command")
	}

	if write == true {
		err = writeIndex(STORAGE_PATH, index)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func readIndex(path string) (*Index, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var index Index
	err = json.Unmarshal(f, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// Writes index as JSON in a file at path.
func writeIndex(path string, index *Index) error {
	raw, err := json.Marshal(index)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, raw, 0644)
}
