package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// NNMainView is the main user interface for nn
type NNMainView struct {
	app           *tview.Application
	cmd           *tview.InputField
	statusBar     *tview.TextView
	listNotebooks *tview.List
	listNotes     *tview.List
	preview       *tview.TextView
	root          *tview.Flex

	model      DataModel
	controller Controller
}

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// DataModel is an interface which provides access to the model
type DataModel interface {
	GetNotebooks() ([]Notebook, error)
	GetNotes() ([]Note, error)
	GetFullPath(notebookIndex, noteIndex int) string
	GetContent(notebookIndex, noteIndex int) string
	GetNewNotePath(topic string) string

	SetSelectedNotebook(notebookIndex int)
	SetSelectedNote(notebookIndex, noteIndex int)
	SetNotebookDirty(notebookIndex int)
}

// Controller is an interface for providing business logic
type Controller interface {
	editFile(fileName string)
}

// Run starts the user interface
func (v *NNMainView) Run() error {
	return v.app.Run()
}

// ---------------------------------------------------------------------------------
// EVENT HANDLER
// ---------------------------------------------------------------------------------
func (v *NNMainView) handleEditNote(index int, mainText, secondaryText string, shortcut rune) {

	v.app.Suspend(func() {
		v.controller.editFile(v.model.GetFullPath(
			v.listNotebooks.GetCurrentItem(),
			v.listNotes.GetCurrentItem()))
	})
}

func (v *NNMainView) handleNotesChanged(index int, mainText, secondaryText string, shortcut rune) {

	if v.listNotes.HasFocus() {
		v.preview.SetText(v.model.GetContent(
			v.listNotebooks.GetCurrentItem(),
			v.listNotes.GetCurrentItem()))
		v.preview.ScrollToBeginning()
	}
}

func (v *NNMainView) handleNotebookChanged(index int, mainText, secondaryText string, shortcut rune) {

	v.model.SetSelectedNotebook(index)
	v.listNotes.Clear()
	notes, err := v.model.GetNotes()
	if err != nil {
		panic("error during fetch list")
	}
	for _, n := range notes {
		v.listNotes.AddItem(n.Name, "", 0, nil)
	}
	v.preview.SetText(v.model.GetContent(v.listNotebooks.GetCurrentItem(), 0))
	v.preview.ScrollToBeginning()
}

// NewTui creates a new TUI
// ---------------------------------------------------------------------------------
func NewTui(d DataModel, c Controller) UIRunner {

	view := NNMainView{}
	view.model = d
	view.controller = c

	notebooks, err := d.GetNotebooks()
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
		if event.Rune() == 'n' {
			view.enableCommand("Topic: ")
		}
		// if event.Rune() == ':' {
		// 	view.app.SetFocus(view.input)
		// 	view.listNotebooks.SetCurrentItem(view.listNotebooks.GetCurrentItem() - 1)
		// }
		return event
	})

	// Notes
	notes, err := d.GetNotes()
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

	view.root = tview.NewFlex().SetDirection(tview.FlexRow)
	view.root.AddItem(main, 0, 1, true)

	view.statusBar = tview.NewTextView().SetText("nn - ready")

	view.cmd = tview.NewInputField().
		SetLabel("Topic: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				view.newNote(view.cmd.GetText())
				view.enableStatusBar()
			} else if key == tcell.KeyEsc {
				view.enableStatusBar()
			}
		})
	view.root.AddItem(view.cmd, 1, 2, false)

	view.app = tview.NewApplication().SetRoot(view.root, true)

	view.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// view.app.Stop()
		}
		return event
	})
	view.enableStatusBar()

	return &view
}

// ---------------------------------------------------------------------------------
// HELPER FUNCTIONS
// ---------------------------------------------------------------------------------

func (v *NNMainView) newNote(topic string) {
	filePath := v.model.GetNewNotePath(topic)
	v.app.Suspend(func() {
		v.controller.editFile(filePath)
	})
	v.model.SetNotebookDirty(v.listNotebooks.GetCurrentItem())
}

func (v *NNMainView) enableStatusBar() {
	v.root.RemoveItem(v.cmd)
	v.root.AddItem(v.statusBar, 1, 2, false)
	v.app.SetFocus(v.listNotebooks)
}

func (v *NNMainView) enableCommand(label string) {
	v.root.RemoveItem(v.statusBar)
	v.cmd.SetLabel(label)
	v.cmd.SetText("")
	v.root.AddItem(v.cmd, 1, 2, true)
	v.app.SetFocus(v.cmd)
}
