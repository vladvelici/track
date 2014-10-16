package main

import (
	"fmt"
	"os"
	"strings"
)

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

// Start working on projects. Returns true if the state has changed.
// Only the first project in the array is going to change.
func start(index *Index, project []string) bool {
	if len(project) != 1 {
		fmt.Println("Can only work at a project at a time.")
		return false
	}
	proj := strings.ToLower(project[0])
	last := index.Status(proj)
	if last == nil {
		// create it.
		//        last :=
	}
	return false
}

// Stop working on a project. Returns true if the state has changed.
func stop(index *Index, project []string) bool { return false }

// Add projects. Returns true if the project wasn't there already.
// Note: Project names are case-insensitive for simplicity.
func add(index *Index, projects []string) bool {
	if len(projects) < 1 {
		fmt.Println("Please use `track add <project-name> [<project-name> ...]`")
		return false
	}

	write := false
	for _, p := range projects {
		p = strings.ToLower(p)
		_, ok := index.Slots[p]
		if !ok {
			index.Slots[p] = make([]*Work, 0)
			write = true
		}
	}

	return write
}

func rm(index *Index, projects []string) bool {
	if len(projects) < 1 {
		fmt.Println("Please use `track rm <project-name> [<project-name> ...]`")
		return false
	}

	if !confirm(fmt.Sprintf("Remove projects %s?", projects)) {
		return false
	}

	write := false
	for _, p := range projects {
		p = strings.ToLower(p)
		_, ok := index.Slots[p]
		write = write || ok
		delete(index.Slots, p)
	}

	return write
}
