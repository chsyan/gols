package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int           // cursor position
	files  []fs.DirEntry // current files
	path   string        // current path
}

func updatePath(m model, path string) model {
	m.path = path
	m.cursor = 0
	m.files = getFiles(path)
	return m
}

func getFiles(path string) []fs.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return files
}

func initialModel() model {
	startPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return model{
		path:  startPath,
		files: getFiles(startPath),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle keypresses
	case tea.KeyMsg:
		switch msg.String() {

		// Exiting
		case "ctrl+c", "q":
			return m, tea.Quit

		// Move cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Move cursor down
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

		// Enter directory
		case "right", "l":
			file := m.files[m.cursor]

			if file.IsDir() {
				m = updatePath(m, filepath.Join(m.path, file.Name()))
			} else {
				// TODO: open file (xdg open?)
				panic("Error: Entering files not implemented yet")
			}

		// Go to parent directory
		case "left", "h":
			parent := filepath.Dir(m.path)
			m = updatePath(m, parent)
		}

	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "Current files \n"

	// Iterate over our choices
	for i, file := range m.files {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, file.Name())
	}

	// The footer
	s += fmt.Sprintf("Current path: %s\n", m.path)
	s += fmt.Sprintf("Selected path: %s\n", m.files[m.cursor].Name())

	files := ""
	for _, file := range m.files {
		files += file.Name() + ", "
	}

	s += fmt.Sprintf("Files: %s\n", files)

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
