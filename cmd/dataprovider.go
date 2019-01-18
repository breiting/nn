package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type dataProvider struct {
	notebooks        []Notebook
	selectedNotebook int
	selectedNote     int
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

	if len(files) > 0 {
		t.selectedNotebook = 0
	} else {
		t.selectedNotebook = -1
	}
	return t.notebooks, nil
}

// GetNotes implements DataProvider interface requires to have
// a selected notebook first
func (t *dataProvider) GetNotes() ([]Note, error) {

	var notes []Note
	if t.selectedNotebook == -1 {
		return notes, fmt.Errorf("no notebook selected")
	}
	// return cached result
	if len(t.notebooks[t.selectedNotebook].Notes) > 0 {
		return t.notebooks[t.selectedNotebook].Notes, nil
	}

	err := filepath.Walk(
		NotesDir+
			string(os.PathSeparator)+
			t.notebooks[t.selectedNotebook].Name, func(path string, info os.FileInfo, e error) error {
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
	t.notebooks[t.selectedNotebook].Notes = notes
	if len(notes) > 0 {
		t.selectedNote = 0
	} else {
		t.selectedNote = -1
	}

	return notes, err
}

func (t *dataProvider) GetContent(notebookIndex, noteIndex int) string {

	b, err := ioutil.ReadFile(t.GetFullPath(notebookIndex, noteIndex))
	if err != nil {
		return "nothing to preview"
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

func (t *dataProvider) getNotebook() (*Notebook, error) {

	if t.selectedNotebook < 0 || t.selectedNotebook > len(t.notebooks)-1 {
		return nil, fmt.Errorf("out of range")
	}
	return &t.notebooks[t.selectedNotebook], nil
}

func (t *dataProvider) SetSelectedNotebook(notebookIndex int) {
	t.selectedNotebook = notebookIndex
}

func (t *dataProvider) SetSelectedNote(notebookIndex, noteIndex int) {
	t.SetSelectedNotebook(notebookIndex)
	t.selectedNote = noteIndex
}

func (t *dataProvider) GetNewNotePath(topic string) string {

	notebook, err := t.getNotebook()
	if err != nil {
		return "dummy.md"
	}
	return NotesDir +
		string(os.PathSeparator) +
		notebook.Name +
		string(os.PathSeparator) +
		"2019-01-11" + // TODO
		"-" +
		topic +
		".md"
}
