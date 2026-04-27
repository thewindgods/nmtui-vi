package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nmtui-vi/internal/nmcli"
	"nmtui-vi/internal/theme"
)

type connItem struct {
	conn nmcli.Connection
}

func (i connItem) Title() string { return i.conn.Name }
func (i connItem) Description() string {
	if i.conn.Device != "" {
		return i.conn.Type + " [" + i.conn.Device + "]"
	}
	return i.conn.Type
}
func (i connItem) FilterValue() string { return i.conn.Name }

type ConnectionsModel struct {
	list          list.Model
	back          bool
	editUUID      string
	newConnType   string
	showOverlay   bool
	overlayCursor int
	err           error
	status        string
	width         int
	height        int
}

var overlayConnTypes = []struct{ label, value string }{
	{"Ethernet", "ethernet"},
	{"WiFi", "wifi"},
}

const connsVerticalMargin = 4

type connsLoadedMsg []nmcli.Connection
type errMsg error
type connActionDoneMsg struct{}

func loadConnections() tea.Cmd {
	return func() tea.Msg {
		conns, err := nmcli.ListConnections()
		if err != nil {
			return errMsg(err)
		}
		return connsLoadedMsg(conns)
	}
}

func connActionCmd(fn func() error) tea.Cmd {
	return func() tea.Msg {
		fn()
		return connActionDoneMsg{}
	}
}

func NewConnectionsModel(w, h int) ConnectionsModel {
	l := newStyledList("Connections", w, max(0, h-connsVerticalMargin))
	return ConnectionsModel{list: l, width: w, height: h}
}

func (m ConnectionsModel) Init() tea.Cmd {
	return loadConnections()
}

func (m ConnectionsModel) Update(msg tea.Msg) (ConnectionsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height-connsVerticalMargin)
		return m, nil

	case connsLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, c := range msg {
			items[i] = connItem{c}
		}
		m.list.SetItems(items)
		m.status = ""

	case errMsg:
		m.err = msg

	case connActionDoneMsg:
		m.status = "Done."
		return m, loadConnections()

	case tea.KeyMsg:
		if m.showOverlay {
			switch msg.String() {
			case "j", "down":
				if m.overlayCursor < len(overlayConnTypes)-1 {
					m.overlayCursor++
				}
			case "k", "up":
				if m.overlayCursor > 0 {
					m.overlayCursor--
				}
			case "enter":
				m.newConnType = overlayConnTypes[m.overlayCursor].value
				m.showOverlay = false
			case "esc", "q":
				m.showOverlay = false
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "esc":
			m.back = true
			return m, nil
		case "enter":
			if sel, ok := m.list.SelectedItem().(connItem); ok {
				m.editUUID = sel.conn.UUID
				return m, nil
			}
		case "a":
			if sel, ok := m.list.SelectedItem().(connItem); ok {
				m.status = "Activating..."
				return m, connActionCmd(func() error { return nmcli.ActivateConnection(sel.conn.UUID) })
			}
		case "d":
			if sel, ok := m.list.SelectedItem().(connItem); ok {
				m.status = "Deactivating..."
				return m, connActionCmd(func() error { return nmcli.DeactivateConnection(sel.conn.UUID) })
			}
		case "D":
			if sel, ok := m.list.SelectedItem().(connItem); ok {
				m.status = "Deleting..."
				return m, connActionCmd(func() error { return nmcli.DeleteConnection(sel.conn.UUID) })
			}
		case "n":
			m.showOverlay = true
			m.overlayCursor = 0
			return m, nil
		case "r":
			m.status = "Refreshing..."
			return m, loadConnections()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ConnectionsModel) Back() bool         { return m.back }
func (m ConnectionsModel) EditUUID() string   { return m.editUUID }
func (m ConnectionsModel) NewConnType() string { return m.newConnType }

func (m ConnectionsModel) View() string {
	if m.err != nil {
		return errorView(m.err)
	}
	status := ""
	if m.status != "" {
		status = "\n  " + theme.StatusStyle.Render(m.status)
	}
	help := "\n" + theme.HelpStyle.Render("  enter: edit  a: activate  d: deactivate  D: delete  n: new  r: refresh  q: back")
	base := "\n" + m.list.View() + status + help

	if m.showOverlay {
		return overlayView(base, m.overlayCursor, m.width, m.height)
	}
	return base
}

// ── overlay styles ────────────────────────────────────────────────────────────

var (
	overlayBorderStyle       lipgloss.Style
	overlayTitleStyle        lipgloss.Style
	overlayItemStyle         lipgloss.Style
	overlayItemSelectedStyle lipgloss.Style
)

func initOverlayStyles() {
	overlayBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Mauve).
		Background(theme.Base).
		Padding(1, 3)

	overlayTitleStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		MarginBottom(1)

	overlayItemStyle = lipgloss.NewStyle().
		Foreground(theme.Text).
		PaddingLeft(2)

	overlayItemSelectedStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		PaddingLeft(1)
}

func overlayView(base string, cursor, w, h int) string {
	var sb strings.Builder
	sb.WriteString(overlayTitleStyle.Render("New Connection Type"))
	sb.WriteString("\n")
	for i, t := range overlayConnTypes {
		if i == cursor {
			sb.WriteString(overlayItemSelectedStyle.Render("> " + t.label))
		} else {
			sb.WriteString(overlayItemStyle.Render(t.label))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(theme.HelpStyle.Render("enter: select  esc: cancel"))

	box := overlayBorderStyle.Render(sb.String())
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, box,
		lipgloss.WithWhitespaceBackground(theme.Base))
}
