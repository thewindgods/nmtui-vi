package theme

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"nmtui-vi/internal/config"
)

var (
	Mauve    lipgloss.Color
	Lavender lipgloss.Color
	Text     lipgloss.Color
	Subtext1 lipgloss.Color
	Subtext0 lipgloss.Color
	Overlay1 lipgloss.Color
	Overlay0 lipgloss.Color
	Surface2 lipgloss.Color
	Surface1 lipgloss.Color
	Surface0 lipgloss.Color
	Base     lipgloss.Color
	Mantle   lipgloss.Color
	Green    lipgloss.Color
	Red      lipgloss.Color
	Peach    lipgloss.Color
)

var (
	TitleStyle        lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	DimItemStyle      lipgloss.Style
	HelpStyle         lipgloss.Style
	StatusStyle       lipgloss.Style
	ErrorStyle        lipgloss.Style
	SuccessStyle      lipgloss.Style
	BorderStyle       lipgloss.Style
	InputLabelStyle   lipgloss.Style
	InputStyle        lipgloss.Style
)

// Init applies the given color config and rebuilds all styles.
// Must be called before starting the bubbletea program.
func Init(c config.Colors) {
	Mauve    = lipgloss.Color(c.Accent)
	Lavender = lipgloss.Color(c.Secondary)
	Text     = lipgloss.Color(c.Text)
	Subtext1 = lipgloss.Color(c.Subtext)
	Subtext0 = lipgloss.Color(c.Muted)
	Overlay1 = lipgloss.Color(c.Subtle)
	Overlay0 = lipgloss.Color(c.Faint)
	Surface2 = lipgloss.Color(c.Border)
	Surface1 = lipgloss.Color(c.Surface)
	Surface0 = lipgloss.Color(c.Panel)
	Base     = lipgloss.Color(c.Background)
	Mantle   = lipgloss.Color(c.Dark)
	Green    = lipgloss.Color(c.Green)
	Red      = lipgloss.Color(c.Red)
	Peach    = lipgloss.Color(c.Orange)

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Base).
		Background(Mauve).
		Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
		Foreground(Text).
		PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
		Foreground(Mauve).
		Bold(true).
		PaddingLeft(1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(Mauve)

	DimItemStyle = lipgloss.NewStyle().
		Foreground(Subtext0).
		PaddingLeft(2)

	HelpStyle = lipgloss.NewStyle().
		Foreground(Overlay1)

	StatusStyle = lipgloss.NewStyle().
		Foreground(Peach)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(Red)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(Green)

	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Surface2).
		Padding(0, 1)

	InputLabelStyle = lipgloss.NewStyle().
		Foreground(Mauve).
		Bold(true)

	InputStyle = lipgloss.NewStyle().
		Foreground(Text)
}

func NewDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.NormalTitle = ItemStyle
	d.Styles.NormalDesc = DimItemStyle
	d.Styles.SelectedTitle = SelectedItemStyle
	d.Styles.SelectedDesc = SelectedItemStyle.Foreground(Lavender).Bold(false)
	d.Styles.DimmedTitle = DimItemStyle
	d.Styles.DimmedDesc = DimItemStyle.Foreground(Overlay0)

	return d
}
