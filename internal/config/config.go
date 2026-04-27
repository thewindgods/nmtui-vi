package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Colors holds all customizable palette colors.
// Defaults are Catppuccin Mocha.
type Colors struct {
	Accent     string
	Secondary  string
	Text       string
	Subtext    string
	Muted      string
	Subtle     string
	Faint      string
	Border     string
	Surface    string
	Panel      string
	Background string
	Dark       string
	Green      string
	Red        string
	Orange     string
}

var defaults = Colors{
	Accent:     "#cba6f7",
	Secondary:  "#b4befe",
	Text:       "#cdd6f4",
	Subtext:    "#bac2de",
	Muted:      "#a6adc8",
	Subtle:     "#7f849c",
	Faint:      "#6c7086",
	Border:     "#585b70",
	Surface:    "#45475a",
	Panel:      "#313244",
	Background: "#1e1e2e",
	Dark:       "#181825",
	Green:      "#a6e3a1",
	Red:        "#f38ba8",
	Orange:     "#fab387",
}

// Load reads ~/.config/nmtui-vi/config and returns the resolved Colors.
// Missing keys fall back to Catppuccin Mocha defaults.
func Load() Colors {
	c := defaults

	home, err := os.UserHomeDir()
	if err != nil {
		return c
	}

	f, err := os.Open(filepath.Join(home, ".config", "nmtui-vi", "config"))
	if err != nil {
		return c
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(parts[0]))
		val := strings.TrimSpace(parts[1])
		switch key {
		case "accent":
			c.Accent = val
		case "secondary":
			c.Secondary = val
		case "text":
			c.Text = val
		case "subtext":
			c.Subtext = val
		case "muted":
			c.Muted = val
		case "subtle":
			c.Subtle = val
		case "faint":
			c.Faint = val
		case "border":
			c.Border = val
		case "surface":
			c.Surface = val
		case "panel":
			c.Panel = val
		case "background":
			c.Background = val
		case "dark":
			c.Dark = val
		case "green":
			c.Green = val
		case "red":
			c.Red = val
		case "orange":
			c.Orange = val
		}
	}

	return c
}
