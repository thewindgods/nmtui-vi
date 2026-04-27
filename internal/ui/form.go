package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nmtui-vi/internal/theme"
)

type fieldKind int

const (
	fieldSection  fieldKind = iota
	fieldInput
	fieldPassword
	fieldSelect
	fieldCheckbox
	fieldButton
)

type formField struct {
	kind       fieldKind
	label      string
	key        string
	value      string
	options    []string
	optLabels  []string
	input      textinput.Model
	dependsOn  string
	showWhen   string
	revealedBy string
	virtual    bool
}

type formModel struct {
	fields      []formField
	cursor      int
	editing     bool
	offset      int
	width       int
	height      int
	lastPressed string
}

// ── constructors ────────────────────────────────────────────────────────────

func sectionField(label string) formField {
	return formField{kind: fieldSection, label: label}
}

func inputField(label, key, value, placeholder string) formField {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.TextStyle = theme.InputStyle
	return formField{kind: fieldInput, label: label, key: key, value: value, input: ti}
}

func inputFieldCond(label, key, value, placeholder, dependsOn, showWhen string) formField {
	f := inputField(label, key, value, placeholder)
	f.dependsOn = dependsOn
	f.showWhen = showWhen
	return f
}

func passwordField(label, key, value, dependsOn, showWhen, revealedBy string) formField {
	ti := textinput.New()
	ti.Placeholder = "password"
	ti.EchoMode = textinput.EchoPassword
	ti.TextStyle = theme.InputStyle
	return formField{kind: fieldPassword, label: label, key: key, value: value, input: ti, dependsOn: dependsOn, showWhen: showWhen, revealedBy: revealedBy}
}

func checkboxField(label, key, dependsOn, showWhen string) formField {
	return formField{kind: fieldCheckbox, label: label, key: key, value: "false", virtual: true, dependsOn: dependsOn, showWhen: showWhen}
}

func buttonField(label, key string) formField {
	return formField{kind: fieldButton, label: label, key: key, virtual: true}
}

func selectField(label, key, value string, options, optLabels []string) formField {
	return formField{kind: fieldSelect, label: label, key: key, value: value, options: options, optLabels: optLabels}
}

func selectFieldCond(label, key, value string, options, optLabels []string, dependsOn, showWhen string) formField {
	f := selectField(label, key, value, options, optLabels)
	f.dependsOn = dependsOn
	f.showWhen = showWhen
	return f
}

func newFormModel(fields []formField, w, h int) formModel {
	inputW := max(10, w-28)
	for i := range fields {
		fields[i].input.Width = inputW
	}
	return formModel{fields: fields, width: w, height: h}
}

// ── querying ─────────────────────────────────────────────────────────────────

func (f *formModel) getValue(key string) string {
	for _, ff := range f.fields {
		if ff.key == key {
			return ff.value
		}
	}
	return ""
}

func (f *formModel) isVisible(i int) bool {
	ff := f.fields[i]
	if ff.dependsOn == "" {
		return true
	}
	return f.getValue(ff.dependsOn) == ff.showWhen
}

func (f *formModel) focusable() []int {
	var out []int
	for i, ff := range f.fields {
		if ff.kind != fieldSection && f.isVisible(i) {
			out = append(out, i)
		}
	}
	return out
}

func (f *formModel) currentIdx() int {
	fv := f.focusable()
	if f.cursor < 0 || f.cursor >= len(fv) {
		return -1
	}
	return fv[f.cursor]
}

func (f *formModel) Values() map[string]string {
	out := make(map[string]string)
	for _, ff := range f.fields {
		if ff.key != "" && ff.kind != fieldSection && !ff.virtual {
			out[ff.key] = ff.value
		}
	}
	return out
}

// ── scrolling ────────────────────────────────────────────────────────────────

func (f *formModel) scrollToCursor(fv []int) {
	if f.cursor >= len(fv) {
		return
	}
	targetFieldIdx := fv[f.cursor]
	row := 0
	for i := range f.fields {
		if !f.isVisible(i) {
			continue
		}
		if i == targetFieldIdx {
			break
		}
		row++
	}
	if row < f.offset {
		f.offset = row
	} else if row >= f.offset+f.height {
		f.offset = row - f.height + 1
	}
}

// ── update ───────────────────────────────────────────────────────────────────

func (f formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		inputW := max(10, msg.Width-28)
		for i := range f.fields {
			f.fields[i].input.Width = inputW
		}
		return f, nil

	case tea.KeyMsg:
		if f.editing {
			switch msg.String() {
			case "esc", "enter":
				idx := f.currentIdx()
				if idx >= 0 {
					f.fields[idx].value = f.fields[idx].input.Value()
					f.fields[idx].input.Blur()
				}
				f.editing = false
				return f, nil
			default:
				idx := f.currentIdx()
				if idx >= 0 {
					var cmd tea.Cmd
					f.fields[idx].input, cmd = f.fields[idx].input.Update(msg)
					return f, cmd
				}
			}
			return f, nil
		}

		switch msg.String() {
		case "j", "down":
			fv := f.focusable()
			if f.cursor < len(fv)-1 {
				f.cursor++
				f.scrollToCursor(fv)
			}
		case "k", "up":
			if f.cursor > 0 {
				fv := f.focusable()
				f.cursor--
				f.scrollToCursor(fv)
			}
		case "enter", "i", " ":
			idx := f.currentIdx()
			if idx < 0 {
				break
			}
			switch f.fields[idx].kind {
			case fieldInput, fieldPassword:
				f.editing = true
				f.fields[idx].input.SetValue(f.fields[idx].value)
				f.fields[idx].input.Width = max(10, f.width-28)
				f.fields[idx].input.Focus()
				f.fields[idx].input.CursorEnd()
				return f, textinput.Blink
			case fieldSelect:
				opts := f.fields[idx].options
				for i, opt := range opts {
					if opt == f.fields[idx].value {
						f.fields[idx].value = opts[(i+1)%len(opts)]
						return f, nil
					}
				}
				if len(opts) > 0 {
					f.fields[idx].value = opts[0]
				}
			case fieldCheckbox:
				if f.fields[idx].value == "true" {
					f.fields[idx].value = "false"
				} else {
					f.fields[idx].value = "true"
				}
			case fieldButton:
				f.lastPressed = f.fields[idx].key
			}
		}

	default:
		if f.editing {
			idx := f.currentIdx()
			if idx >= 0 {
				var cmd tea.Cmd
				f.fields[idx].input, cmd = f.fields[idx].input.Update(msg)
				return f, cmd
			}
		}
	}
	return f, nil
}

// ── styles ───────────────────────────────────────────────────────────────────

var (
	formSectionStyle          lipgloss.Style
	formLabelStyle            lipgloss.Style
	formLabelSelectedStyle    lipgloss.Style
	formValueStyle            lipgloss.Style
	formValueSelectedStyle    lipgloss.Style
	formCursorStyle           lipgloss.Style
	formButtonStyle           lipgloss.Style
	formButtonUnselectedStyle lipgloss.Style
)

func initFormStyles() {
	formSectionStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		PaddingLeft(2)

	formLabelStyle = lipgloss.NewStyle().
		Foreground(theme.Subtext1).
		Width(20)

	formLabelSelectedStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		Width(20)

	formValueStyle = lipgloss.NewStyle().
		Foreground(theme.Text)

	formValueSelectedStyle = lipgloss.NewStyle().
		Foreground(theme.Lavender).
		Bold(true)

	formCursorStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true)

	formButtonStyle = lipgloss.NewStyle().
		Foreground(theme.Base).
		Background(theme.Mauve).
		Bold(true).
		Padding(0, 1)

	formButtonUnselectedStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, true).
		BorderForeground(theme.Mauve).
		Foreground(theme.Mauve).
		Padding(0, 1)
}

// ── rendering ────────────────────────────────────────────────────────────────

func (f formModel) View() string {
	fv := f.focusable()
	cursorFieldIdx := -1
	if f.cursor < len(fv) {
		cursorFieldIdx = fv[f.cursor]
	}

	var lines []string
	for i, ff := range f.fields {
		if !f.isVisible(i) {
			continue
		}
		if ff.kind == fieldSection {
			sep := strings.Repeat("─", max(0, f.width-len(ff.label)-6))
			lines = append(lines, formSectionStyle.Render("── "+ff.label+" "+sep))
		} else {
			lines = append(lines, f.renderField(i, i == cursorFieldIdx))
		}
	}

	start := f.offset
	if start > len(lines) {
		start = len(lines)
	}
	end := min(start+f.height, len(lines))

	result := make([]string, f.height)
	copy(result, lines[start:end])
	return strings.Join(result, "\n")
}

func (f formModel) renderField(idx int, selected bool) string {
	ff := f.fields[idx]

	if ff.kind == fieldButton {
		cur := "  "
		if selected {
			cur = formCursorStyle.Render("> ")
			return cur + formButtonStyle.Render(" "+ff.label+" ")
		}
		return cur + formButtonUnselectedStyle.Render(ff.label)
	}

	cur := "  "
	if selected {
		cur = formCursorStyle.Render("> ")
	}

	labelStyle := formLabelStyle
	if selected {
		labelStyle = formLabelSelectedStyle
	}
	label := labelStyle.Render(ff.label + ":")

	var value string
	if f.editing && selected && (ff.kind == fieldInput || ff.kind == fieldPassword) {
		if ff.kind == fieldPassword {
			inputCopy := ff.input
			if ff.revealedBy != "" && f.getValue(ff.revealedBy) == "true" {
				inputCopy.EchoMode = textinput.EchoNormal
			} else {
				inputCopy.EchoMode = textinput.EchoPassword
			}
			value = inputCopy.View()
		} else {
			value = ff.input.View()
		}
	} else {
		raw := f.displayValue(idx)
		if selected {
			value = formValueSelectedStyle.Render(raw)
		} else {
			value = formValueStyle.Render(raw)
		}
	}

	return cur + label + "  " + value
}

func (f formModel) displayValue(idx int) string {
	ff := f.fields[idx]
	switch ff.kind {
	case fieldSelect:
		for i, opt := range ff.options {
			if opt == ff.value {
				if i < len(ff.optLabels) {
					return ff.optLabels[i]
				}
				return opt
			}
		}
		if len(ff.optLabels) > 0 {
			return theme.HelpStyle.Render(ff.optLabels[0])
		}
		return theme.HelpStyle.Render("(select)")
	case fieldPassword:
		if ff.value == "" {
			return theme.HelpStyle.Render("(empty)")
		}
		if ff.revealedBy != "" && f.getValue(ff.revealedBy) == "true" {
			return ff.value
		}
		return strings.Repeat("●", min(len(ff.value), 12))
	case fieldCheckbox:
		if ff.value == "true" {
			return theme.SuccessStyle.Render("[x]")
		}
		return theme.HelpStyle.Render("[ ]")
	default:
		if ff.value == "" {
			return theme.HelpStyle.Render("(empty)")
		}
		return ff.value
	}
}
