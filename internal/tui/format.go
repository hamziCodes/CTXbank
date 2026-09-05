package tui

import (
	"os"
)

// IsTTY returns true if stdout is connected to a terminal/console.
func IsTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// BoxChars holds boundary characters for terminal card rendering.
type BoxChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	Bullet      string
	Check       string
	Cross       string
}

var (
	// UnicodeBox is used when running on a modern interactive terminal.
	UnicodeBox = BoxChars{
		TopLeft:     "┌─",
		TopRight:    "─┐",
		BottomLeft:  "└─",
		BottomRight: "─┘",
		Horizontal:  "─",
		Vertical:    "│",
		Bullet:      "●",
		Check:       "✓",
		Cross:       "✗",
	}

	// ASCIIBox is used when output is redirected to a pipe, CI runner, or agent.
	ASCIIBox = BoxChars{
		TopLeft:     "+-",
		TopRight:    "-+",
		BottomLeft:  "+-",
		BottomRight: "-+",
		Horizontal:  "-",
		Vertical:    "|",
		Bullet:      "*",
		Check:       "+",
		Cross:       "x",
	}
)

// GetBox returns UnicodeBox if interactive TTY, otherwise ASCIIBox.
func GetBox() BoxChars {
	if IsTTY() {
		return UnicodeBox
	}
	return ASCIIBox
}
