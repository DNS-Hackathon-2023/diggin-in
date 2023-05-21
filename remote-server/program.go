package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type Program struct {
	filename string
}

// Save program to file
func (p *Program) Save(text []byte) error {
	// Write program text to file
	if err := os.WriteFile(p.filename, text, 0644); err != nil {
		return err
	}
	return nil
}

// Load program from file
func (p *Program) Load() (string, error) {
	// Read program text from file
	text, err := os.ReadFile(p.filename)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

// Calculate the sha256 hash of the program as
// the 'ID'.
func (p *Program) ID() (string, error) {
	// Read program text from file
	text, err := p.Load()
	if err != nil {
		return "", err
	}
	// Calculate sha256 sum
	sum := sha256.Sum256([]byte(text))
	id := fmt.Sprintf("%x", sum)

	return id, nil
}

// DefaultProgram returns a new Program with default filename
func DefaultProgram() *Program {
	return &Program{filename: "program.txt"}
}
