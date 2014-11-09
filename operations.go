package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	StatusWorking      = iota
	StatusNotWorking   = iota
	StatusDoesNotExist = iota
)

// Project names are case sensitive in the Index.
type Index struct {
	Slots map[string]*Project
}

// Get a list of all projects.
func (i *Index) Projects() []string {
	projects := make([]string, 0, len(i.Slots))
	for k, _ := range i.Slots {
		projects = append(projects, k)
	}
	return projects
}

func (i *Index) ProjectExists(project string) bool {
	_, ok := i.Slots[project]
	return ok
}

// Get the status of a project. Returns one of the constats Status*
func (i *Index) Status(project string) int {
	proj, ok := i.Slots[project]
	if !ok {
		return StatusDoesNotExist
	}
	return proj.Status()
}

func (i *Index) CurrentProject() (*Project, string) {
	for name, proj := range i.Slots {
		if proj.Status() == StatusWorking {
			return proj, name
		}
	}
	return nil, ""
}

// Print current status on terminal
func (i *Index) PrintStatus() {
	for name, proj := range i.Slots {
		if proj.Status() == StatusWorking {
			fmt.Print(" -> ")
		} else {
			fmt.Print("    ")
		}
		fmt.Printf("%s\n", name)
	}
}

// Get the last work entry from a project. Returns nil if it has none
// or the project doesn't exist.
func (i *Index) LastWorkEntry(project string) *Work {
	proj, ok := i.Slots[project]
	if !ok {
		return nil
	}
	return proj.LastWorkEntry()
}

func (i *Index) StartWorking(project string) (bool, error) {
	proj, ok := i.Slots[project]
	if !ok {
		return false, fmt.Errorf("Project %s does not exist. Add it with `track add %s`.", project, project)
	}
	return proj.Work()
}

func (i *Index) StopWorking() (bool, error) {
	proj, _ := i.CurrentProject()
	if proj == nil {
		return false, fmt.Errorf("Not currently working on any project.")
	}
	return proj.Stop()
}

// Add projects. Returns true if the project wasn't there already.
// Note: Project names are made case-insensitive on this method
func (index *Index) AddProjects(projects []string) (bool, error) {
	if len(projects) < 1 {
		return false, fmt.Errorf("Please use `track add <project-name> [<project-name> ...]`")
	}

	write := false
	for _, p := range projects {
		p = strings.ToLower(p)
		_, ok := index.Slots[p]
		if !ok {
			tmp := make(Project, 0)
			index.Slots[p] = &tmp
			write = true
		}
	}

	return write, nil
}

// Remove a project from the index.
func (index *Index) RemoveProjects(projects []string) (bool, error) {
	if len(projects) < 1 {
		return false, fmt.Errorf("Please use `track rm <project-name> [<project-name> ...]`")
	}

	write := false
	for _, p := range projects {
		p = strings.ToLower(p)
		_, ok := index.Slots[p]
		write = write || ok
		delete(index.Slots, p)
	}

	return write, nil
}

// Type Project represents a project with all its work entries.
type Project []*Work

// Returns the last work entry or nil if there's none.
func (p *Project) LastWorkEntry() *Work {
	if len(*p) == 0 {
		return nil
	}
	return (*p)[len(*p)-1]
}

func (p *Project) Add(w *Work) {
	*p = append(*p, w)
}

func (p *Project) Status() int {
	w := p.LastWorkEntry()
	if w == nil || !w.End.IsZero() {
		return StatusNotWorking
	}
	return StatusWorking
}

func (p *Project) Work() (bool, error) {
	if p.Status() == StatusWorking {
		return false, fmt.Errorf("Already working since %s", p.LastWorkEntry().Start.Local())
	}
	p.Add(&Work{Start: time.Now().UTC()})
	return true, nil
}

func (p *Project) Stop() (bool, error) {
	if p.Status() == StatusNotWorking {
		return false, fmt.Errorf("Not working.")
	}
	w := p.LastWorkEntry()
	w.End = time.Now().UTC()
	return true, nil
}

type Work struct {
	Start time.Time
	End   time.Time
}

// Creates an empty index and writes the index file.
func createIndex(path string) (*Index, error) {
	index := &Index{
		Slots: make(map[string]*Project),
	}

	err := writeIndex(path, index)
	return index, err
}

// Delete the index. Ask the user before calling this function.
func deleteIndex(path string) error {
	return os.Remove(path)
}
