package main

// Notebook defines the overall notebook structure
type Notebook struct {
	Name  string
	Notes []Note
	Dirty bool
}
