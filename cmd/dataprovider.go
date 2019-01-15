package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type dataProvider struct {
	notebooks []Notebook
}

// NewDataProvider creates a new data provider
func NewDataProvider() DataProvider {
	return &dataProvider{}
}

func acceptedPath(f os.FileInfo) bool {
	return f.IsDir() &&
		!strings.Contains(f.Name(), ".git")
}

// GetNotebooks implements DataProvider interface
func (t *dataProvider) GetNotebooks() ([]Notebook, error) {

	files, err := ioutil.ReadDir(NotesDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if acceptedPath(f) {
			t.notebooks = append(t.notebooks, Notebook{Name: f.Name()})
		}
	}

	return t.notebooks, nil
}

// GetNotes implements DataProvider interface
func (t *dataProvider) GetNotes(notebookIndex int) ([]Note, error) {

	var notes []Note

	err := filepath.Walk(NotesDir+string(os.PathSeparator)+t.notebooks[notebookIndex].Name, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir)
		if info.Mode().IsRegular() {
			notes = append(notes, Note{Name: info.Name()})
		}
		return nil
	})

	// cache result
	t.notebooks[notebookIndex].Notes = notes
	return notes, err
}

func (t *dataProvider) GetContent(notebookIndex, noteIndex int) string {

	b, err := ioutil.ReadFile(t.GetFullPath(notebookIndex, noteIndex))
	if err != nil {
		return "error opening file"
	}

	return string(b)
}

func (t *dataProvider) GetFullPath(notebookIndex, noteIndex int) string {
	if notebookIndex < 0 || notebookIndex > len(t.notebooks)-1 {
		return ""
	}
	notebook := t.notebooks[notebookIndex]

	if noteIndex < 0 || noteIndex > len(notebook.Notes)-1 {
		return ""
	}

	return NotesDir +
		string(os.PathSeparator) +
		notebook.Name +
		string(os.PathSeparator) +
		notebook.Notes[noteIndex].Name
}
