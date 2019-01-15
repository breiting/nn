package main

import (
	"strconv"
)

type testData struct{}

// NewTestDataProvider creates a new test data provider
func NewTestDataProvider() DataProvider {
	return &testData{}
}

// GetNotebooks implements DataProvider interface
func (t *testData) GetNotebooks() ([]Notebook, error) {

	var notebooks []Notebook

	for i := 0; i < 10; i++ {
		name := "Book " + strconv.FormatInt((int64)(i), 10)
		notebooks = append(notebooks, Notebook{Name: name})
	}

	return notebooks, nil
}

// GetNotes implements DataProvider interface
func (t *testData) GetNotes(notebookIndex int) ([]Note, error) {

	var notes []Note

	for i := 0; i < 100; i++ {
		name := "Note " + strconv.FormatInt((int64)(i), 10) + " of Book " + strconv.FormatInt((int64)(notebookIndex), 10)
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
