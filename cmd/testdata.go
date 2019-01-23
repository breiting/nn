package main

import (
	"strconv"
)

type testData struct{}

// NewTestDataModel creates a new test data provider
func NewTestDataModel() DataModel {
	return &testData{}
}

// GetNotebooks implements DataModel interface
func (t *testData) GetNotebooks() ([]Notebook, error) {

	var notebooks []Notebook

	for i := 0; i < 10; i++ {
		name := "Book " + strconv.FormatInt((int64)(i), 10)
		notebooks = append(notebooks, Notebook{Name: name})
	}

	return notebooks, nil
}

// GetNotes implements DataModel interface
func (t *testData) GetNotes() ([]Note, error) {

	var notes []Note

	for i := 0; i < 100; i++ {
		name := "Note " + strconv.FormatInt((int64)(i), 10)
		notes = append(notes, Note{Name: name})
	}

	return notes, nil
}

func (t *testData) GetContent(nbIndex, nIndex int) string {
	return "Test content"
}

func (t *testData) GetFullPath(notebookIndex, noteIndex int) string {
	return "test.txt"
}

func (t *testData) SetSelectedNotebook(notebookIndex int) {
}

func (t *testData) SetSelectedNote(notebookIndex, noteIndex int) {
}

func (t *testData) GetNewNotePath(topic string) string {
	return "test.txt"
}

func (t *testData) SetNotebookDirty(notebookIndex int) {
}
