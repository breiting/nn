package main

import (
	"fmt"
	tui "github.com/marcusolsson/tui-go"
	"os"
	"os/exec"
)

// NNMainView is the main user interface for nn
type NNMainView struct {
	ui            tui.UI
	root          *tui.Box
	topic         *tui.TextEdit
	statusBar     *tui.StatusBar
	listNotebooks *tui.List
	listNotes     *tui.List
	preview       *tui.TextEdit
	focusWidgets  []tui.Widget
	data          DataProvider
}

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// DataProvider is an interface which provides access to the data
type DataProvider interface {
	GetNotebooks() ([]Notebook, error)
	GetNotes() ([]Note, error)
	GetFullPath(notebookIndex, noteIndex int) string
	GetContent(notebookIndex, noteIndex int) string
	GetNewNotePath(topic string) string

	SetSelectedNotebook(notebookIndex int)
	SetSelectedNote(notebookIndex, noteIndex int)
}

// Run starts the user interface
func (v *NNMainView) Run() error {
	return v.ui.Run()
}

// NewTui creates a new TUI
func NewTui(c DataProvider) UIRunner {

	view := NNMainView{}
	view.data = c

	notebooks, err := c.GetNotebooks()
	if err != nil {
		panic("no notebooks found")
	}

	// Notebooks
	view.listNotebooks = tui.NewList()
	for _, n := range notebooks {
		view.listNotebooks.AddItems(n.Name)
	}
	view.focusWidgets = append(view.focusWidgets, view.listNotebooks)
	view.listNotebooks.SetFocused(true)

	sidebar := tui.NewVBox(view.listNotebooks)
	sidebar.SetBorder(true)
	sidebar.SetSizePolicy(tui.Minimum, tui.Expanding)

	// Notes
	notes, err := c.GetNotes()
	view.listNotes = tui.NewList()
	for _, n := range notes {
		view.listNotes.AddItems(n.Name)
	}
	view.focusWidgets = append(view.focusWidgets, view.listNotes)

	notesBox := tui.NewVBox(view.listNotes)
	notesBox.SetBorder(true)

	// Preview
	view.preview = tui.NewTextEdit()
	view.preview.SetSizePolicy(tui.Expanding, tui.Expanding)
	view.preview.SetText("preview")
	view.preview.SetWordWrap(true)
	previewBox := tui.NewVBox(view.preview)
	previewBox.SetBorder(true)
	previewBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	tableBox := tui.NewHBox(sidebar, notesBox, previewBox)
	tableBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	// Statusbar
	view.statusBar = tui.NewStatusBar("")
	view.statusBar.SetText("[nn]")
	view.statusBar.SetPermanentText("[press esc or q to quit]")

	view.root = tui.NewVBox(
		tableBox,
		view.statusBar,
	)

	th := tui.NewTheme()
	th.SetStyle("table.cell.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	th.SetStyle("list.item", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	th.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})

	view.ui, _ = tui.New(view.root)
	view.ui.SetTheme(th)
	view.setKeybindingsDefault()

	// nbList.OnItemActivated(func(l *tui.List) {
	// })
	view.listNotebooks.OnSelectionChanged(func(l *tui.List) {
		view.data.SetSelectedNotebook(l.Selected())
		view.listNotes.RemoveItems()
		notes, err := c.GetNotes()
		if err != nil {
			panic("error during fetch list")
		}
		for _, n := range notes {
			view.listNotes.AddItems(n.Name)
		}
	})

	view.listNotes.OnSelectionChanged(func(l *tui.List) {
		view.preview.SetText(c.GetContent(
			view.listNotebooks.Selected(),
			view.listNotes.Selected()))
	})

	view.listNotebooks.Select(0)
	return &view
}

func (v *NNMainView) unFocusAll() {
	for _, w := range v.focusWidgets {
		w.SetFocused(false)
	}
}

func (v *NNMainView) setKeybindingsDefault() {
	v.ui.ClearKeybindings()
	v.ui.SetKeybinding("Esc", func() { v.ui.Quit() })
	v.ui.SetKeybinding("q", func() { v.ui.Quit() })
	v.ui.SetKeybinding("h", func() {
		v.unFocusAll()
		v.listNotebooks.SetFocused(true)
	})
	v.ui.SetKeybinding("Enter", func() {
		if v.listNotes.IsFocused() {
			v.openFile(v.data.GetFullPath(
				v.listNotebooks.Selected(),
				v.listNotes.Selected()))
		}
	})

	v.ui.SetKeybinding("l", func() {
		v.unFocusAll()
		v.listNotes.SetFocused(true)
		v.listNotes.Select(0)
	})

	v.ui.SetKeybinding("n", func() {
		if v.listNotebooks.IsFocused() {
			label := tui.NewLabel("Topic: ")
			v.topic = tui.NewTextEdit()
			command := tui.NewHBox(label, v.topic, tui.NewSpacer())

			v.setKeybindingsCommand()
			v.unFocusAll()
			v.root.Remove(1)

			v.root.Append(command)
			v.topic.SetFocused(true)
		}
	})
}

func (v *NNMainView) setKeybindingsCommand() {
	v.ui.ClearKeybindings()

	v.ui.SetKeybinding("Esc", func() {
		v.unFocusAll()
		v.root.Remove(1)
		v.root.Append(v.statusBar)
		v.setKeybindingsDefault()
		v.listNotebooks.SetFocused(true)
	})

	v.ui.SetKeybinding("Enter", func() {
		v.newNote(v.topic.Text())
	})
}

func (v *NNMainView) newNote(topic string) {
	filePath := v.data.GetNewNotePath(topic)
	v.openFile(filePath)
}

func (v *NNMainView) openFile(fileName string) {
	var cmd *exec.Cmd
	cmd = exec.Command(Editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	v.ui.Quit()
	err := cmd.Run()
	if err != nil {
		fmt.Println("Couldn't open the file:", err)
		os.Exit(1)
	}
	// TODO sync git
	os.Exit(0)
}
