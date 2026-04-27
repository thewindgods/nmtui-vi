package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"nmtui-vi/internal/theme"
)

// Init rebuilds all package-level styles from the current theme values.
// Must be called after theme.Init().
func Init() {
	initFormStyles()
	initOverlayStyles()
}

func errorView(err error) string {
	return fmt.Sprintf("\n%s\n\nPress q to go back.", theme.ErrorStyle.Render("Error: "+err.Error()))
}

func newStyledList(title string, w, h int) list.Model {
	l := list.New([]list.Item{}, theme.NewDelegate(), w, max(0, h))
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = theme.TitleStyle
	l.KeyMap.Quit.SetKeys("q", "esc")
	return l
}
