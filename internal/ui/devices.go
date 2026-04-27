package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"nmtui-vi/internal/nmcli"
	"nmtui-vi/internal/theme"
)

type devItem struct {
	dev nmcli.Device
}

func (i devItem) Title() string { return i.dev.Name }
func (i devItem) Description() string {
	conn := i.dev.Connection
	if conn == "" {
		conn = "disconnected"
	}
	return fmt.Sprintf("%s — %s — %s", i.dev.Type, i.dev.State, conn)
}
func (i devItem) FilterValue() string { return i.dev.Name }

type DevicesModel struct {
	list list.Model
	back bool
	err  error
}

const devsVerticalMargin = 4

type devsLoadedMsg []nmcli.Device

func loadDevices() tea.Cmd {
	return func() tea.Msg {
		devs, err := nmcli.ListDevices()
		if err != nil {
			return errMsg(err)
		}
		return devsLoadedMsg(devs)
	}
}

func NewDevicesModel(w, h int) DevicesModel {
	l := newStyledList("Device Status", w, max(0, h-devsVerticalMargin))
	return DevicesModel{list: l}
}

func (m DevicesModel) Init() tea.Cmd {
	return loadDevices()
}

func (m DevicesModel) Update(msg tea.Msg) (DevicesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-devsVerticalMargin)
		return m, nil

	case devsLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, d := range msg {
			items[i] = devItem{d}
		}
		m.list.SetItems(items)

	case errMsg:
		m.err = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.back = true
			return m, nil
		case "r":
			return m, loadDevices()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m DevicesModel) Back() bool { return m.back }

func (m DevicesModel) View() string {
	if m.err != nil {
		return errorView(m.err)
	}
	return "\n" + m.list.View() + "\n" + theme.HelpStyle.Render("  r: refresh  q: back")
}
