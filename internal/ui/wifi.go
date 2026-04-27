package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"nmtui-vi/internal/nmcli"
	"nmtui-vi/internal/theme"
)

type wifiItem struct {
	net nmcli.WifiNetwork
}

func (i wifiItem) Title() string {
	if i.net.InUse {
		return theme.SuccessStyle.Render("* " + i.net.SSID)
	}
	return i.net.SSID
}
func (i wifiItem) Description() string {
	return "signal: " + i.net.Signal + "  security: " + i.net.Security
}
func (i wifiItem) FilterValue() string { return i.net.SSID }

type WifiModel struct {
	list        list.Model
	passwordBox textinput.Model
	connecting  bool
	selectedNet string
	back        bool
	err         error
	status      string
}

const wifiVerticalMargin = 4

type wifiLoadedMsg []nmcli.WifiNetwork
type wifiConnectedMsg struct{ err error }

func scanWifi() tea.Cmd {
	return func() tea.Msg {
		nets, err := nmcli.ScanWifi()
		if err != nil {
			return errMsg(err)
		}
		return wifiLoadedMsg(nets)
	}
}

func connectWifi(ssid, password string) tea.Cmd {
	return func() tea.Msg {
		err := nmcli.AddWifiConnection(ssid, password)
		return wifiConnectedMsg{err}
	}
}

func NewWifiModel(w, h int) WifiModel {
	l := newStyledList("WiFi Networks", w, max(0, h-wifiVerticalMargin))

	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.EchoMode = textinput.EchoPassword
	ti.Width = 40
	ti.PromptStyle = theme.InputLabelStyle
	ti.TextStyle = theme.InputStyle

	return WifiModel{list: l, passwordBox: ti}
}

func (m WifiModel) Init() tea.Cmd {
	return scanWifi()
}

func (m WifiModel) Update(msg tea.Msg) (WifiModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-wifiVerticalMargin)
		return m, nil

	case wifiLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, n := range msg {
			items[i] = wifiItem{n}
		}
		m.list.SetItems(items)
		m.status = ""

	case wifiConnectedMsg:
		m.connecting = false
		if msg.err != nil {
			m.status = "error:" + msg.err.Error()
		} else {
			m.status = "connected"
			return m, scanWifi()
		}

	case errMsg:
		m.err = msg

	case tea.KeyMsg:
		if m.connecting {
			switch msg.String() {
			case "enter":
				ssid := m.selectedNet
				pass := m.passwordBox.Value()
				m.passwordBox.Reset()
				m.status = "connecting"
				return m, connectWifi(ssid, pass)
			case "esc":
				m.connecting = false
				m.passwordBox.Reset()
			default:
				var cmd tea.Cmd
				m.passwordBox, cmd = m.passwordBox.Update(msg)
				return m, cmd
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "esc":
			m.back = true
			return m, nil
		case "enter":
			if sel, ok := m.list.SelectedItem().(wifiItem); ok {
				if sel.net.Security != "" && sel.net.Security != "--" {
					m.connecting = true
					m.selectedNet = sel.net.SSID
					m.passwordBox.Focus()
					return m, textinput.Blink
				}
				m.status = "connecting"
				return m, connectWifi(sel.net.SSID, "")
			}
		case "r":
			m.status = "scanning"
			return m, scanWifi()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m WifiModel) Back() bool { return m.back }

func (m WifiModel) View() string {
	if m.err != nil {
		return errorView(m.err)
	}
	if m.connecting {
		return "\n  " + theme.InputLabelStyle.Render("Connect to:") + " " + m.selectedNet +
			"\n\n  " + m.passwordBox.View() +
			"\n\n" + theme.HelpStyle.Render("  enter: confirm  esc: cancel")
	}
	status := ""
	switch m.status {
	case "connecting", "scanning":
		status = "\n  " + theme.StatusStyle.Render(m.status+"...")
	case "connected":
		status = "\n  " + theme.SuccessStyle.Render("Connected!")
	default:
		if m.status != "" {
			status = "\n  " + theme.ErrorStyle.Render(m.status)
		}
	}
	return "\n" + m.list.View() + status + "\n" + theme.HelpStyle.Render("  enter: connect  r: rescan  q: back")
}
