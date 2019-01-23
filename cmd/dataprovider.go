package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type dataModel struct {
	notebooks        []Notebook
	selectedNotebook int
	selectedNote     int
}

// ByDate sorts the notes by date
type ByDate []Note

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Modified.After(a[j].Modified) }

// NewDataModel creates a new data provider
func NewDataModel() DataModel {
	return &dataModel{}
}

func acceptedPath(f os.FileInfo) bool {
	return f.IsDir() &&
		!strings.Contains(f.Name(), ".git") &&
		!strings.Contains(f.Name(), ".template")
}

// GetNotebooks implements DataModel interface
func (t *dataModel) GetNotebooks() ([]Notebook, error) {

	files, err := ioutil.ReadDir(NotesDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if acceptedPath(f) {
			t.notebooks = append(t.notebooks, Notebook{Name: f.Name(), Dirty: true})
		}
	}

	if len(files) > 0 {
		t.selectedNotebook = 0
	} else {
		t.selectedNotebook = -1
	}
	return t.notebooks, nil
}

// GetNotes implements DataModel interface requires to have
// a selected notebook first
func (t *dataModel) GetNotes() ([]Note, error) {

	var notes []Note
	if t.selectedNotebook == -1 {
		return notes, fmt.Errorf("no notebook selected")
	}
	// return cached result
	if len(t.notebooks[t.selectedNotebook].Notes) > 0 && !t.notebooks[t.selectedNotebook].Dirty {
		return t.notebooks[t.selectedNotebook].Notes, nil
	}

	err := filepath.Walk(
		NotesDir+
			string(os.PathSeparator)+
			t.notebooks[t.selectedNotebook].Name, func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			// check if it is a regular file (not dir) and ignore hidden files
			if info.Mode().IsRegular() && !strings.HasPrefix(info.Name(), ".") {
				notes = append(notes, Note{Name: info.Name(), Modified: info.ModTime()})
			}
			return nil
		})

	// sort result
	sort.Sort(ByDate(notes))

	// cache result
	t.notebooks[t.selectedNotebook].Notes = notes
	t.notebooks[t.selectedNotebook].Dirty = false

	if len(notes) > 0 {
		t.selectedNote = 0
	} else {
		t.selectedNote = -1
	}

	return notes, err
}

func (t *dataModel) GetContent(notebookIndex, noteIndex int) string {

	b, err := ioutil.ReadFile(t.GetFullPath(notebookIndex, noteIndex))
	if err != nil {
		return "nothing to preview"
	}

	return string(b)
}

func (t *dataModel) GetFullPath(notebookIndex, noteIndex int) string {
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

func (t *dataModel) getNotebook() (*Notebook, error) {

	if t.selectedNotebook < 0 || t.selectedNotebook > len(t.notebooks)-1 {
		return nil, fmt.Errorf("out of range")
	}
	return &t.notebooks[t.selectedNotebook], nil
}

func (t *dataModel) SetSelectedNotebook(notebookIndex int) {
	t.selectedNotebook = notebookIndex
}

func (t *dataModel) SetSelectedNote(notebookIndex, noteIndex int) {
	t.SetSelectedNotebook(notebookIndex)
	t.selectedNote = noteIndex
}

func (t *dataModel) GetNewNotePath(topic string) string {

	notebook, err := t.getNotebook()
	if err != nil {
		return "dummy.md"
	}
	now := time.Now()
	dateStr := fmt.Sprintf("%4d-%02d-%02d", now.Year(), now.Month(), now.Day())
	return NotesDir +
		string(os.PathSeparator) +
		notebook.Name +
		string(os.PathSeparator) +
		dateStr +
		"-" +
		topic +
		".md"
}

func (t *dataModel) SetNotebookDirty(notebookIndex int) {
	t.notebooks[notebookIndex].Dirty = true
}
