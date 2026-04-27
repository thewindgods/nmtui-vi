package nmcli

import (
	"fmt"
	"os/exec"
	"strings"
)

type Connection struct {
	Name   string
	UUID   string
	Type   string
	Device string
}

type Device struct {
	Name       string
	Type       string
	State      string
	Connection string
}

func ListConnections() ([]Connection, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "NAME,UUID,TYPE,DEVICE", "connection", "show").Output()
	if err != nil {
		return nil, err
	}
	var conns []Connection
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 {
			continue
		}
		conns = append(conns, Connection{
			Name:   parts[0],
			UUID:   parts[1],
			Type:   parts[2],
			Device: parts[3],
		})
	}
	return conns, nil
}

func ListDevices() ([]Device, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device", "status").Output()
	if err != nil {
		return nil, err
	}
	var devs []Device
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 {
			continue
		}
		devs = append(devs, Device{
			Name:       parts[0],
			Type:       parts[1],
			State:      parts[2],
			Connection: parts[3],
		})
	}
	return devs, nil
}

func ActivateConnection(uuid string) error {
	return exec.Command("nmcli", "connection", "up", uuid).Run()
}

func DeactivateConnection(uuid string) error {
	return exec.Command("nmcli", "connection", "down", uuid).Run()
}

func DeleteConnection(uuid string) error {
	return exec.Command("nmcli", "connection", "delete", uuid).Run()
}

func AddWifiConnection(ssid, password string) error {
	args := []string{"device", "wifi", "connect", ssid}
	if password != "" {
		args = append(args, "password", password)
	}
	return exec.Command("nmcli", args...).Run()
}

func ScanWifi() ([]WifiNetwork, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "SSID,SIGNAL,SECURITY,IN-USE", "device", "wifi", "list").Output()
	if err != nil {
		return nil, err
	}
	var nets []WifiNetwork
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 || parts[0] == "" {
			continue
		}
		nets = append(nets, WifiNetwork{
			SSID:     parts[0],
			Signal:   parts[1],
			Security: parts[2],
			InUse:    parts[3] == "*",
		})
	}
	return nets, nil
}

type WifiNetwork struct {
	SSID     string
	Signal   string
	Security string
	InUse    bool
}

func GetConnectionDetails(uuid string) (map[string]string, error) {
	out, err := exec.Command("nmcli", "--show-secrets", "-t", "-f", "all", "connection", "show", uuid).Output()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := line[:idx]
		val := strings.ReplaceAll(line[idx+1:], `\:`, `:`)
		result[key] = strings.TrimSpace(val)
	}
	return result, nil
}

func ModifyConnection(uuid string, settings map[string]string) error {
	args := []string{"connection", "modify", uuid}
	for k, v := range settings {
		if k != "" {
			args = append(args, k, v)
		}
	}
	out, err := exec.Command("nmcli", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func AddConnection(connType, name string, settings map[string]string) error {
	args := []string{"connection", "add", "type", connType, "con-name", name}
	for k, v := range settings {
		if k != "" && v != "" {
			args = append(args, k, v)
		}
	}
	out, err := exec.Command("nmcli", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}
