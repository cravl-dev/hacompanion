package main

import (
	"fmt"
	"hacompanion/entity"
	"log"
	"net"
)

// Config contains all values from the configuration file.
type Config struct {
	HomeAssistant homeassistantConfig            `toml:"homeassistant"`
	Companion     companionConfig                `toml:"companion"`
	Notifications notificationsConfig            `toml:"notifications"`
	Sensors       map[string]entity.SensorConfig `toml:"sensor"`
	Script        map[string]entity.ScriptConfig `toml:"script"`
}

type homeassistantConfig struct {
	DeviceName string `toml:"device_name"`
	Token      string `toml:"token"`
	Host       string `toml:"host"`
}

type companionConfig struct {
	UpdateInterval   duration `toml:"update_interval"`
	RegistrationFile homePath `toml:"registration_file"`
}

type notificationsConfig struct {
	Listen  string `toml:"listen"`
	PushURL string `toml:"push_url"`
}

func getLocalIp() (string, error) {
	// Credit: https://gist.github.com/jniltinho/9787946?permalink_comment_id=2243615#gistcomment-2243615
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetPushUrl returns the pushUrl if set in the config or tries to guess it
func (c Config) GetPushUrl() (string, error) {
	// Use whatever is set in the config file if set
	if c.Notifications.PushURL != "" {
		return c.Notifications.PushURL, nil
	}
	// Get listen port
	port := ":8234"
	if c.Notifications.Listen != "" {
		port = c.Notifications.Listen
	}

	localIp, err := getLocalIp()
	if err != nil {
		log.Println("failed to determine local IP. Please set notifications.push_url in your config")
		return "", err
	}

	return fmt.Sprintf("http://%s%s/notifications", localIp, port), nil
}
