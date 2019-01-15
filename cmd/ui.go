package main

import (
	"fmt"
	tui "github.com/marcusolsson/tui-go"
	"os"
	"os/exec"
)

var ui tui.UI

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// DataProvider is an interface which provides access to the data
type DataProvider interface {
	GetNotebooks() ([]Notebook, error)
	GetNotes(notebookIndex int) ([]Note, error)
	GetFullPath(notebookIndex, noteIndex int) string
	GetContent(notebookIndex, noteIndex int) string
}

// NewTui creates a new TUI
func NewTui(c DataProvider) tui.UI {

	notebooks, err := c.GetNotebooks()
	if err != nil {
		panic("no notebooks found")
	}

	// Notebooks
	nbList := tui.NewList()
	for _, n := range notebooks {
		nbList.AddItems(n.Name)
	}
	nbList.SetFocused(true)

	sidebar := tui.NewVBox(nbList)
	sidebar.SetBorder(true)
	sidebar.SetSizePolicy(tui.Minimum, tui.Expanding)

	// Notes
	notes, err := c.GetNotes(0)
	noteList := tui.NewList()
	for _, n := range notes {
		noteList.AddItems(n.Name)
	}
	notesView := tui.NewVBox(noteList)
	notesView.SetBorder(true)

	// Preview
	buffer := tui.NewTextEdit()
	buffer.SetSizePolicy(tui.Expanding, tui.Expanding)
	buffer.SetText("preview")
	buffer.SetWordWrap(true)
	bufferView := tui.NewVBox(buffer)
	bufferView.SetBorder(true)
	bufferView.SetSizePolicy(tui.Expanding, tui.Expanding)

	tableBox := tui.NewHBox(sidebar, notesView, bufferView)

	// Statusbar
	status := tui.NewStatusBar("")
	status.SetText("[nn]")
	status.SetPermanentText("[press esc or q to quit]")

	root := tui.NewVBox(
		tableBox,
		status,
	)

	th := tui.NewTheme()
	th.SetStyle("table.cell.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	th.SetStyle("list.item", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	th.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})

	ui, _ = tui.New(root)
	ui.SetTheme(th)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("r", func() { ui.Repaint() })
	ui.SetKeybinding("h", func() {
		nbList.SetFocused(true)
		noteList.SetFocused(false)
	})
	ui.SetKeybinding("Enter", func() {
		if noteList.IsFocused() {
			openFile(c.GetFullPath(nbList.Selected(), noteList.Selected()))
		}
	})

	ui.SetKeybinding("l", func() {
		nbList.SetFocused(false)
		noteList.SetFocused(true)
		noteList.Select(0)
	})

	nbList.OnItemActivated(func(l *tui.List) {
	})
	nbList.OnSelectionChanged(func(l *tui.List) {
		noteList.RemoveItems()
		notes, err := c.GetNotes(l.Selected())
		if err != nil {
			panic("error during fetch list")
		}
		for _, n := range notes {
			noteList.AddItems(n.Name)
		}

	})

	noteList.OnItemActivated(func(l *tui.List) {
	})
	noteList.OnSelectionChanged(func(l *tui.List) {
		buffer.SetText(c.GetContent(nbList.Selected(), l.Selected()))
	})
	nbList.Select(0)

	return ui
}

func openFile(fileName string) {
	var cmd *exec.Cmd
	cmd = exec.Command(Editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ui.Quit()
	err := cmd.Run()
	if err != nil {
		fmt.Println("Couldn't open the file:", err)
		os.Exit(1)
	}
	// TODO sync git
	os.Exit(0)
}
