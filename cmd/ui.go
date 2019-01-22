package main

import (
	"fmt"
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"os"
	"os/exec"
)

// NNMainView is the main user interface for nn
type NNMainView struct {
	app           *tview.Application
	cmd           *tview.InputField
	statusBar     *tview.TextView
	data          DataProvider
	listNotebooks *tview.List
	listNotes     *tview.List
	preview       *tview.TextView
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
	return v.app.Run()
}

func (v *NNMainView) handleEditNote(index int, mainText, secondaryText string, shortcut rune) {

	v.app.Suspend(func() {
		v.openFile(v.data.GetFullPath(
			v.listNotebooks.GetCurrentItem(),
			v.listNotes.GetCurrentItem()))
	})
}

func (v *NNMainView) handleNotesChanged(index int, mainText, secondaryText string, shortcut rune) {

	if v.listNotes.HasFocus() {
		v.preview.SetText(v.data.GetContent(
			v.listNotebooks.GetCurrentItem(),
			v.listNotes.GetCurrentItem()))
		v.preview.ScrollToBeginning()
	}
}

func (v *NNMainView) handleNotebookChanged(index int, mainText, secondaryText string, shortcut rune) {

	v.data.SetSelectedNotebook(index)
	v.listNotes.Clear()
	notes, err := v.data.GetNotes()
	if err != nil {
		panic("error during fetch list")
	}
	for _, n := range notes {
		v.listNotes.AddItem(n.Name, "", 0, nil)
	}
	v.preview.SetText(v.data.GetContent(v.listNotebooks.GetCurrentItem(), 0))
	v.preview.ScrollToBeginning()
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
	view.listNotebooks = tview.NewList().ShowSecondaryText(false)
	for _, n := range notebooks {
		view.listNotebooks.AddItem(n.Name, "", 0, nil)
	}
	view.listNotebooks.SetChangedFunc(view.handleNotebookChanged)
	view.listNotebooks.SetBorder(true)
	view.listNotebooks.SetCurrentItem(0)

	view.listNotebooks.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			view.app.SetFocus(view.listNotebooks)
		}
		if event.Rune() == 'l' {
			view.app.SetFocus(view.listNotes)
			view.listNotes.SetCurrentItem(0)
		}
		if event.Rune() == 'q' {
			view.app.Stop()
		}
		// if event.Rune() == ':' {
		// 	view.app.SetFocus(view.input)
		// 	view.listNotebooks.SetCurrentItem(view.listNotebooks.GetCurrentItem() - 1)
		// }
		return event
	})

	// Notes
	notes, err := c.GetNotes()
	view.listNotes = tview.NewList().ShowSecondaryText(false)
	for _, n := range notes {
		view.listNotes.AddItem(n.Name, "", 0, nil)
	}
	view.listNotes.SetBorder(true)
	view.listNotes.SetChangedFunc(view.handleNotesChanged)
	view.listNotes.SetSelectedFunc(view.handleEditNote)

	view.listNotes.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			view.app.SetFocus(view.listNotebooks)
		} else if event.Rune() == 'q' {
			view.app.Stop()
		}
		return event
	})

	// Preview
	view.preview = tview.NewTextView().SetText("nothing to show").SetWrap(true)
	view.preview.SetBorder(true)

	main := tview.NewFlex()
	main.AddItem(view.listNotebooks, 0, 1, true)
	main.AddItem(view.listNotes, 0, 2, false)
	main.AddItem(view.preview, 0, 3, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(main, 0, 1, true)

	view.statusBar = tview.NewTextView().SetText("nn - ready")

	view.cmd = tview.NewInputField().
		SetLabel("Topic: ").
		SetDoneFunc(func(key tcell.Key) {
			view.app.Stop()
		})

	root.AddItem(view.statusBar, 1, 2, false)

	view.app = tview.NewApplication().SetRoot(root, true)

	view.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			view.app.Stop()
		}
		return event
	})

	//
	// // Preview
	// view.preview = tui.NewTextEdit()
	// view.preview.SetSizePolicy(tui.Expanding, tui.Expanding)
	// view.preview.SetText("preview")
	// view.preview.SetWordWrap(true)
	// previewBox := tui.NewVBox(view.preview)
	// previewBox.SetBorder(true)
	// previewBox.SetSizePolicy(tui.Expanding, tui.Expanding)
	//
	// tableBox := tui.NewHBox(sidebar, notesBox, previewBox)
	// tableBox.SetSizePolicy(tui.Expanding, tui.Expanding)
	//
	// // Statusbar
	// view.statusBar = tui.NewStatusBar("")
	// view.statusBar.SetText("[nn]")
	// view.statusBar.SetPermanentText("[press esc or q to quit]")
	//
	// view.root = tui.NewVBox(
	// 	tableBox,
	// 	view.statusBar,
	// )
	//
	// th := tui.NewTheme()
	// th.SetStyle("table.cell.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	// th.SetStyle("list.item", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	// th.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	//
	// view.ui, _ = tui.New(view.root)
	// view.ui.SetTheme(th)
	// view.setKeybindingsDefault()
	//
	// // nbList.OnItemActivated(func(l *tui.List) {
	// // })
	// view.listNotebooks.OnSelectionChanged(func(l *tui.List) {
	// 	view.data.SetSelectedNotebook(l.Selected())
	// 	view.listNotes.RemoveItems()
	// 	notes, err := c.GetNotes()
	// 	if err != nil {
	// 		panic("error during fetch list")
	// 	}
	// 	for _, n := range notes {
	// 		view.listNotes.AddItems(n.Name)
	// 	}
	// })
	//
	// view.listNotes.OnSelectionChanged(func(l *tui.List) {
	// 	view.preview.SetText(c.GetContent(
	// 		view.listNotebooks.Selected(),
	// 		view.listNotes.Selected()))
	// })
	//
	// view.listNotebooks.Select(0)
	return &view
}

func (v *NNMainView) openFile(fileName string) {
	var cmd *exec.Cmd
	cmd = exec.Command(Editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Couldn't open the file:", err)
		os.Exit(1)
	}
}
