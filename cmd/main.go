package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	// DefaultEditor sets the default editor, can be overwritten by EDITOR
	DefaultEditor = "vim"

	// DefaultNotesDir holds the note data (default is ~/notes), can be overwritten by NNDIR
	DefaultNotesDir = "notes"

	// CmdTree is the Linux command for listing a directory
	CmdTree = "tree"
)

var (
	// Editor stores the editor which should be used
	Editor = DefaultEditor
	// NotesDir stores the directory where all the notes are stored
	NotesDir = DefaultNotesDir
)

func help() {

	fmt.Println(`nn - simple note taking

	nn         interactive tui

	nn init    create new note environment (not implemented)
	nn check   checks if everything is setup properly (not implemented)

	nn new     create new note (not implemented)
	nn sync    sync note with git server (not implemented)

	nn show    shows the content of the notes directory

	nn help    prints this message
	`)
}

func init() {

	editorEnv := os.Getenv("EDITOR")
	notesDirEnv := os.Getenv("NNDIR")
	user, _ := user.Current()

	if editorEnv != "" {
		Editor = editorEnv
	}
	if notesDirEnv != "" {
		NotesDir = notesDirEnv
	} else {
		NotesDir = filepath.Join(user.HomeDir, DefaultNotesDir)
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
		show()
	case "help":
		help()
	default:
		help()
	}
}

func show() {
	var cmd *exec.Cmd
	cmd = exec.Command(CmdTree, NotesDir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("Cannot execute command:", err)
		os.Exit(1)
	}

}

func printConfig() {
	fmt.Printf("Configuration\n\n")
	fmt.Println("	EDITOR           " + Editor)
	fmt.Println("	Notes directory  " + NotesDir)
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
		// dataModel := NewTestModel()
		dataModel := NewDataModel()
		controller := NewController()
		NewTui(dataModel, controller).Run()
	} else if len(os.Args) == 2 {
		parseAction(os.Args[1])
	}
}
