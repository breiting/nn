package main

// Tui is a terminal user interface
type Tui interface {
	Start()
}

type tui struct{}

// NewTui creates a new TUI
func NewTui() Tui {
	return &tui{}
}

func (t *tui) Start() {
	NotImplemented()
}
