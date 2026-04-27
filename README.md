# nmtui-vi

A terminal UI for NetworkManager with vi keybindings, written in Go.
Uses `nmcli` as its backend.

<img width="967" height="577" alt="2026-04-27-195502_hyprshot" src="https://github.com/user-attachments/assets/8f3575c4-fdb4-436b-a037-cfe6a62a06df" />

## What it does

- Browse, create, edit, and delete network connections
- Activate and deactivate connections
- View device status (name, type, state, active connection)
- Scan for and connect to WiFi networks
- Edit IPv4 and IPv6 settings (method, address, gateway, DNS)
- Edit WiFi security settings (WPA/WPA2 Personal and Enterprise)
- Theming via a config file

## What it does not do

- VPN connections (OpenVPN, WireGuard, IPsec)
- Bond, bridge, VLAN, or team connection types
- Proxy settings
- Setting the system hostname

## Requirements

- Go 1.21 or later
- NetworkManager with `nmcli` available in `$PATH`

## Installation

```
git clone <repository>
cd nmtui-vi
go build -o nmtui-vi .
```

Move the binary somewhere on your `$PATH`:

```
mv nmtui-vi ~/.local/bin/
```

## Usage

```
nmtui-vi
```

## Navigation

| Key | Action |
|-----|--------|
| `j` / `down` | Move down |
| `k` / `up` | Move up |
| `enter` | Select / confirm |
| `q` / `esc` | Go back |

### Connections screen

| Key | Action |
|-----|--------|
| `enter` | Edit selected connection |
| `a` | Activate selected connection |
| `d` | Deactivate selected connection |
| `D` | Delete selected connection |
| `n` | New connection |
| `r` | Refresh list |

### Connection editor

| Key | Action |
|-----|--------|
| `i` / `enter` | Edit field / cycle option |
| `space` | Cycle select option |
| `S` | Save connection |
| `a` | Activate connection |
| `esc` | Cancel edit / go back |

## Theming

Colors are read from `~/.config/nmtui-vi/config` at startup. The file uses a
simple `key = value` format. Lines beginning with `#` are comments. Any
omitted key falls back to the Catppuccin Mocha default.

Example — Catppuccin Macchiato:

```
accent     = #c6a0f6
secondary  = #b7bdf8
text       = #cad3f5
subtext    = #b8c0e0
muted      = #a5adcb
subtle     = #8087a2
faint      = #6e738d
border     = #5b6078
surface    = #494d64
panel      = #363a4f
background = #24273a
dark       = #1e2030
green      = #a6da95
red        = #ed8796
orange     = #f5a97f
```

### Color keys

| Key | Role |
|-----|------|
| `accent` | Primary accent — titles, cursor, selected items |
| `secondary` | Secondary accent — selected field values |
| `text` | Body text |
| `subtext` | Labels and secondary text |
| `muted` | Dimmed list items |
| `subtle` | Help bar text |
| `faint` | Dimmed descriptions |
| `border` | Borders and separators |
| `surface` | Elevated surface background |
| `panel` | Panel background |
| `background` | Main background |
| `dark` | Deep background layer |
| `green` | Success messages |
| `red` | Error messages |
| `orange` | Status and loading messages |
