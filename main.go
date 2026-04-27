package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"nmtui-vi/internal/config"
	"nmtui-vi/internal/theme"
	"nmtui-vi/internal/ui"
)

type appState int

const (
	stateMain appState = iota
	stateConnections
	stateDevices
	stateWifi
	stateEditor
)

type model struct {
	state       appState
	width       int
	height      int
	mainMenu    ui.MainMenuModel
	connections ui.ConnectionsModel
	devices     ui.DevicesModel
	wifi        ui.WifiModel
	editor      ui.EditorModel
}

func newModel() model {
	return model{
		state:    stateMain,
		mainMenu: ui.NewMainMenu(0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return m.mainMenu.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sz, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = sz.Width, sz.Height
	}

	switch m.state {
	case stateMain:
		var cmd tea.Cmd
		m.mainMenu, cmd = m.mainMenu.Update(msg)
		if m.mainMenu.Done() {
			switch m.mainMenu.Chosen() {
			case ui.ScreenConnections:
				m.state = stateConnections
				m.connections = ui.NewConnectionsModel(m.width, m.height)
				return m, m.connections.Init()
			case ui.ScreenDevices:
				m.state = stateDevices
				m.devices = ui.NewDevicesModel(m.width, m.height)
				return m, m.devices.Init()
			case ui.ScreenWifi:
				m.state = stateWifi
				m.wifi = ui.NewWifiModel(m.width, m.height)
				return m, m.wifi.Init()
			default:
				return m, tea.Quit
			}
		}
		return m, cmd

	case stateConnections:
		var cmd tea.Cmd
		m.connections, cmd = m.connections.Update(msg)
		if m.connections.Back() {
			return m.goToMain()
		}
		if uuid := m.connections.EditUUID(); uuid != "" {
			m.state = stateEditor
			m.editor = ui.NewEditorModel(uuid, "", m.width, m.height)
			return m, m.editor.Init()
		}
		if t := m.connections.NewConnType(); t != "" {
			m.state = stateEditor
			m.editor = ui.NewEditorModel("", t, m.width, m.height)
			return m, m.editor.Init()
		}
		return m, cmd

	case stateDevices:
		var cmd tea.Cmd
		m.devices, cmd = m.devices.Update(msg)
		if m.devices.Back() {
			return m.goToMain()
		}
		return m, cmd

	case stateWifi:
		var cmd tea.Cmd
		m.wifi, cmd = m.wifi.Update(msg)
		if m.wifi.Back() {
			return m.goToMain()
		}
		return m, cmd

	case stateEditor:
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		if m.editor.Back() {
			m.state = stateConnections
			m.connections = ui.NewConnectionsModel(m.width, m.height)
			return m, m.connections.Init()
		}
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateMain:
		return m.mainMenu.View()
	case stateConnections:
		return m.connections.View()
	case stateDevices:
		return m.devices.View()
	case stateWifi:
		return m.wifi.View()
	case stateEditor:
		return m.editor.View()
	}
	return ""
}

func (m model) goToMain() (model, tea.Cmd) {
	m.state = stateMain
	m.mainMenu = ui.NewMainMenu(m.width, m.height)
	return m, m.mainMenu.Init()
}

func main() {
	theme.Init(config.Load())
	ui.Init()
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
