package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type Screen int

const (
	ScreenMain Screen = iota
	ScreenConnections
	ScreenDevices
	ScreenWifi
)

const mainMenuVerticalMargin = 2

type MainMenuModel struct {
	list   list.Model
	chosen Screen
	done   bool
}

func NewMainMenu(w, h int) MainMenuModel {
	items := []list.Item{
		menuItem{"Edit Connections", "Add, edit, or remove network connections"},
		menuItem{"Device Status", "View device states"},
		menuItem{"Scan WiFi", "List available wireless networks"},
	}

	l := newStyledList("NetworkManager TUI", w, max(0, h-mainMenuVerticalMargin))
	l.SetItems(items)
	l.KeyMap.Quit.SetKeys("q", "ctrl+c")

	return MainMenuModel{list: l}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (MainMenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-mainMenuVerticalMargin)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(menuItem); ok {
				switch i.title {
				case "Edit Connections":
					m.chosen = ScreenConnections
				case "Device Status":
					m.chosen = ScreenDevices
				case "Scan WiFi":
					m.chosen = ScreenWifi
				}
				m.done = true
			}
		case "q", "ctrl+c":
			m.done = true
			m.chosen = -1
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) Done() bool     { return m.done }
func (m MainMenuModel) Chosen() Screen { return m.chosen }

func (m MainMenuModel) View() string {
	return "\n" + m.list.View()
}
