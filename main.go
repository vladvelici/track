package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	STORAGE_PATH = "./data.json"
)

type Index struct {
	Slots map[string][]*Work
}

// Get a list of all projects.
func (i *Index) Projects() []string {
	projects := make([]string, 0, len(i.Slots))
	for k, _ := range i.Slots {
		projects = append(projects, k)
	}
	return projects
}

// Get a list of work slots for project. Same as Index.Slots[<project>]
func (i *Index) ProjectSlots(project string) (res []*Work, ok bool) {
	res, ok = i.Slots[project]
	return res, ok
}

type Work struct {
	Start time.Time
	End   time.Time
}

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
	case "rm":
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
		write = start(index, os.Args[1:])
	case "stop":
		write = stop(index, os.Args[1:])
	case "add":
		write = add(index, os.Args[1:])
	case "rm":
		write = rm(index, os.Args[1:])
	case "status":
		status(index, os.Args[1:])
	default:
		fmt.Println("Unkown command")
	}

	if write == true {
		writeIndex(STORAGE_PATH, index)
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

// Creates an empty index and writes the index file.
func createIndex(path string) (*Index, error) {
	index := &Index{
		Slots: make(map[string][]*Work),
	}

	err := writeIndex(path, index)
	return index, err
}

// Delete the index. Ask the user before calling this function.
func deleteIndex(path string) error {
	return os.Remove(path)
}

// Print the current status of the system on the screen.
func status(index *Index, projects []string) {}

// Start working on a project. Returns true if the state has changed.
func start(index *Index, project []string) bool { return false }

// Stop working on a project. Returns true if the state has changed.
func stop(index *Index, project []string) bool { return false }

// Add projects. Returns true if the project wasn't there already.
// Note: Project names are case-insensitive for simplicity.
func add(index *Index, projects []string) bool {
	if len(projects) < 1 {
		fmt.Println("Please use `track add <project-name> [<project-name> ...]`")
		return false
	}

	return false
}

func rm(index *Index, projects []string) bool {
	if len(projects) < 1 {
		fmt.Println("Please use `track rm <project-name> [<project-name> ...]`")
		return false
	}

	if !confirm(fmt.Sprintf("Remove projects %s?", projects)) {
		return false
	}

	for _, p := range projects {
		_, ok := index.Slots[p]
		if ok {
			return true
		}
	}

	return false
}
