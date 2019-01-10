package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	// DefaultEditor sets the default editor, can be overwritten by EDITOR
	DefaultEditor = "vim"

	// DefaultNotesDir holds the note data (default is ~/notes), can be overwritten by NNDIR
	DefaultNotesDir = "notes"
)

var (
	editor   = DefaultEditor
	notesDir = DefaultNotesDir
)

func help() {

	fmt.Println(`nn - simple note taking

	nn         interactive tui

	nn init    create new note environment
	nn check   checks if everything is setup properly

	nn new     create new note
	nn sync    sync note with git server

	nn show    shows the content of the notes directory

	nn help    prints this message
	`)
}

func init() {

	editorEnv := os.Getenv("EDITOR")
	notesDirEnv := os.Getenv("NNDIR")
	user, _ := user.Current()

	if editorEnv != "" {
		editor = editorEnv
	}
	if notesDirEnv != "" {
		notesDir = notesDirEnv
	} else {
		notesDir = filepath.Join(user.HomeDir, DefaultNotesDir)
	}
}

func parseAction(action string) {
	switch action {
	case "init":
		fmt.Println("nn init - not implemented")
		break
	case "check":
		fmt.Println("nn check - not implemented")
		break
	case "new":
		fmt.Println("nn new - not implemented")
		break
	case "sync":
		fmt.Println("nn sync - not implemented")
		break
	case "show":
		fmt.Println("nn show - not implemented")
		break
	case "help":
	default:
		help()
		break
	}
}

func printConfig() {
	fmt.Printf("Configuration\n\n")
	fmt.Println("	EDITOR           " + editor)
	fmt.Println("	Notes directory  " + notesDir)
}

// NotImplemented exits the application
func NotImplemented() {
	fmt.Printf("Not implemented!\n\n")
	help()
	printConfig()
	os.Exit(0)
}

func main() {

	if len(os.Args) == 1 {
		tui := NewTui()
		tui.Start()
	} else if len(os.Args) == 2 {
		parseAction(os.Args[1])
	}
}
