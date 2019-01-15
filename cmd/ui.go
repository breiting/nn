package main

import (
	tui "github.com/marcusolsson/tui-go"
	"io/ioutil"
)

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

type notebook struct {
	name  string
	count int64
}

// Open returns an UIRunner from a given directory
func Open(path string) (UIRunner, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var notebooks []notebook
	for _, f := range files {
		notebooks = append(notebooks, notebook{name: f.Name()})
	}
	return newTuiUI(notebooks), nil
}

func newTuiUI(notebooks []notebook) tui.UI {

	var str []string
	for _, n := range notebooks {
		str = append(str, n.name)
		// t.AppendRow(
		// 	tui.NewLabel(n.name),
		// 	tui.NewLabel(strconv.FormatInt(n.count, 10)),
		// )
	}
	l := tui.NewList()
	l.AddItems(str...)
	l.SetFocused(true)

	status := tui.NewStatusBar("")
	status.SetText("[press enter to switch to selected branch]")
	status.SetPermanentText("[press esc or q to quit]")
	tableBox := tui.NewVBox(l, tui.NewSpacer())
	tableBox.SetBorder(true)
	root := tui.NewVBox(
		l,
		status,
	)

	th := tui.NewTheme()
	th.SetStyle("table.cell.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	th.SetStyle("list.item", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	th.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})

	ui, _ := tui.New(root)
	ui.SetTheme(th)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	l.OnItemActivated(func(l *tui.List) {
	})
	l.OnSelectionChanged(func(l *tui.List) {
	})
	l.Select(0)

	return ui
}
