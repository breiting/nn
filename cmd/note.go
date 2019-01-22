package main

import (
	"time"
)

// Note defines the actual note data structure
type Note struct {
	Name     string
	FilePath string
	Dirty    bool
	Modified time.Time
}
