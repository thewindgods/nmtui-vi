package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"nmtui-vi/internal/nmcli"
	"nmtui-vi/internal/theme"
)

const editorVerticalMargin = 5

type EditorModel struct {
	form     formModel
	uuid     string
	connType string
	loading  bool
	back     bool
	status   string
	width    int
	height   int
}

type detailsLoadedMsg struct {
	details  map[string]string
	connType string
}
type editorResultMsg struct {
	err     error
	success string
}

func NewEditorModel(uuid, connType string, w, h int) EditorModel {
	em := EditorModel{uuid: uuid, connType: connType, width: w, height: h}
	if uuid == "" {
		em.form = buildForm(connType, nil, w, h-editorVerticalMargin)
	} else {
		em.loading = true
	}
	return em
}

func (m EditorModel) Init() tea.Cmd {
	if m.uuid == "" {
		return nil
	}
	uuid := m.uuid
	return func() tea.Msg {
		details, err := nmcli.GetConnectionDetails(uuid)
		if err != nil {
			return errMsg(err)
		}
		return detailsLoadedMsg{details: details, connType: details["connection.type"]}
	}
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)
		m.form.height = m.height - editorVerticalMargin
		return m, cmd

	case detailsLoadedMsg:
		m.loading = false
		m.connType = msg.connType
		m.form = buildForm(msg.connType, msg.details, m.width, m.height-editorVerticalMargin)
		return m, nil

	case editorResultMsg:
		if msg.err != nil {
			m.status = theme.ErrorStyle.Render("Error: " + msg.err.Error())
		} else {
			m.status = theme.SuccessStyle.Render(msg.success)
		}
		return m, nil

	case errMsg:
		m.loading = false
		m.status = theme.ErrorStyle.Render("Error: " + msg.Error())
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			if msg.String() == "esc" || msg.String() == "q" {
				m.back = true
			}
			return m, nil
		}
		switch msg.String() {
		case "esc", "q":
			if m.form.editing {
				var cmd tea.Cmd
				m.form, cmd = m.form.Update(msg)
				return m, cmd
			}
			m.back = true
			return m, nil
		case "S":
			if !m.form.editing {
				return m, m.saveCmd()
			}
		case "a":
			if !m.form.editing && m.uuid != "" {
				return m, m.activateCmd()
			}
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)
		if m.form.lastPressed == "_save" {
			m.form.lastPressed = ""
			return m, m.saveCmd()
		}
		return m, cmd
	}
	return m, nil
}

func (m EditorModel) activateCmd() tea.Cmd {
	uuid := m.uuid
	return func() tea.Msg {
		return editorResultMsg{err: nmcli.ActivateConnection(uuid), success: "Activated"}
	}
}

func (m EditorModel) saveCmd() tea.Cmd {
	values := m.form.Values()
	uuid := m.uuid
	connType := m.connType
	name := values["connection.id"]

	return func() tea.Msg {
		var err error
		if uuid != "" {
			err = nmcli.ModifyConnection(uuid, values)
		} else {
			settings := make(map[string]string)
			for k, v := range values {
				if k != "connection.id" {
					settings[k] = v
				}
			}
			err = nmcli.AddConnection(connType, name, settings)
		}
		return editorResultMsg{err: err, success: "Saved"}
	}
}

func (m EditorModel) Back() bool { return m.back }

func (m EditorModel) View() string {
	title := theme.TitleStyle.Width(m.width).Render(editorTitle(m.uuid, m.connType))

	if m.loading {
		return fmt.Sprintf("\n%s\n\n  %s\n", title, theme.StatusStyle.Render("Loading..."))
	}

	status := ""
	if m.status != "" {
		status = "\n  " + m.status
	}

	helpText := "  j/k: move  i/enter: edit field  enter/space: cycle option  S: save  q/esc: cancel"
	if m.uuid != "" {
		helpText = "  j/k: move  i/enter: edit field  enter/space: cycle option  S: save  a: activate  q/esc: cancel"
	}
	help := "\n" + theme.HelpStyle.Render(helpText)
	return fmt.Sprintf("\n%s\n%s%s%s", title, m.form.View(), status, help)
}

func editorTitle(uuid, connType string) string {
	if uuid == "" {
		return " New " + connTypeLabel(connType) + " Connection "
	}
	return " Edit Connection "
}

func connTypeLabel(t string) string {
	switch t {
	case "802-11-wireless", "wifi":
		return "WiFi"
	case "802-3-ethernet", "ethernet":
		return "Ethernet"
	default:
		if t != "" {
			return strings.ToTitle(t[:1]) + t[1:]
		}
		return "Connection"
	}
}

// ── field builders ────────────────────────────────────────────────────────────

func get(d map[string]string, key, def string) string {
	if d == nil {
		return def
	}
	if v, ok := d[key]; ok && v != "" {
		return v
	}
	return def
}

func buildForm(connType string, d map[string]string, w, h int) formModel {
	switch connType {
	case "802-11-wireless", "wifi":
		return newFormModel(wifiFields(d), w, h)
	default:
		return newFormModel(ethernetFields(d), w, h)
	}
}

func commonFields(d map[string]string) []formField {
	return []formField{
		sectionField("Connection"),
		inputField("Name", "connection.id", get(d, "connection.id", "New Connection"), "connection name"),
		selectField("Autoconnect", "connection.autoconnect",
			get(d, "connection.autoconnect", "yes"),
			[]string{"yes", "no"},
			[]string{"Yes", "No"},
		),
	}
}

func ipv4Fields(d map[string]string) []formField {
	method := get(d, "ipv4.method", "auto")
	return []formField{
		sectionField("IPv4"),
		selectField("Method", "ipv4.method", method,
			[]string{"auto", "manual", "disabled"},
			[]string{"Automatic", "Manual", "Disabled"},
		),
		inputFieldCond("Address/Prefix", "ipv4.addresses",
			get(d, "ipv4.addresses", ""), "e.g. 192.168.1.10/24",
			"ipv4.method", "manual"),
		inputFieldCond("Gateway", "ipv4.gateway",
			get(d, "ipv4.gateway", ""), "e.g. 192.168.1.1",
			"ipv4.method", "manual"),
		inputFieldCond("DNS Servers", "ipv4.dns",
			get(d, "ipv4.dns", ""), "e.g. 8.8.8.8,1.1.1.1",
			"ipv4.method", "manual"),
	}
}

func ipv6Fields(d map[string]string) []formField {
	method := get(d, "ipv6.method", "auto")
	return []formField{
		sectionField("IPv6"),
		selectField("Method", "ipv6.method", method,
			[]string{"auto", "manual", "disabled", "ignore"},
			[]string{"Automatic", "Manual", "Disabled", "Ignore"},
		),
		inputFieldCond("Address/Prefix", "ipv6.addresses",
			get(d, "ipv6.addresses", ""), "e.g. ::1/128",
			"ipv6.method", "manual"),
		inputFieldCond("Gateway", "ipv6.gateway",
			get(d, "ipv6.gateway", ""), "",
			"ipv6.method", "manual"),
		inputFieldCond("DNS Servers", "ipv6.dns",
			get(d, "ipv6.dns", ""), "",
			"ipv6.method", "manual"),
	}
}

func ethernetFields(d map[string]string) []formField {
	fields := commonFields(d)
	fields = append(fields,
		sectionField("Ethernet"),
		inputField("MAC Address", "802-3-ethernet.cloned-mac-address",
			get(d, "802-3-ethernet.cloned-mac-address", ""), "leave blank for default"),
		inputField("MTU", "802-3-ethernet.mtu",
			get(d, "802-3-ethernet.mtu", ""), "leave blank for default"),
	)
	fields = append(fields, ipv4Fields(d)...)
	fields = append(fields, ipv6Fields(d)...)
	return fields
}

func wifiFields(d map[string]string) []formField {
	keyMgmt := get(d, "802-11-wireless-security.key-mgmt", "")

	fields := commonFields(d)
	fields = append(fields,
		sectionField("WiFi"),
		inputField("SSID", "802-11-wireless.ssid",
			get(d, "802-11-wireless.ssid", ""), "network name"),
		selectField("Mode", "802-11-wireless.mode",
			get(d, "802-11-wireless.mode", "infrastructure"),
			[]string{"infrastructure", "ap"},
			[]string{"Client (Infrastructure)", "Access Point"},
		),
		sectionField("Security"),
		selectField("Type", "802-11-wireless-security.key-mgmt",
			keyMgmt,
			[]string{"", "wpa-psk", "wpa-eap"},
			[]string{"None", "WPA/WPA2 Personal", "WPA/WPA2 Enterprise"},
		),
		passwordField("Password", "802-11-wireless-security.psk",
			get(d, "802-11-wireless-security.psk", ""),
			"802-11-wireless-security.key-mgmt", "wpa-psk",
			"_show_password",
		),
		checkboxField("Show password", "_show_password",
			"802-11-wireless-security.key-mgmt", "wpa-psk",
		),
	)
	fields = append(fields, ipv4Fields(d)...)
	fields = append(fields, ipv6Fields(d)...)
	return fields
}
