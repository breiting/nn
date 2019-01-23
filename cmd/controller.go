package main

import (
	"fmt"
	"os"
	"os/exec"
)

type controller struct {
}

// NewController creates a new controller
func NewController() Controller {
	return &controller{}
}

func (c *controller) editFile(fileName string) {
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
