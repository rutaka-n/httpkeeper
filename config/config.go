package config

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
)

type (
	// Config high-level config object
	Config struct {
		Server            Server    `json:"server"`
		InvalidatedTokens []string  `json:"invalidated_tokens"`
		Services          []Service `json:"services"`
	}

	// Server represents configuration of the reverse proxy server itself
	Server struct {
		Addr            string `json:"addr"`
		LogFile         string `json:"log_file"`
		Secret          []byte `json:"secret"`
		Name            string `json:"name"`
		ShutdownTimeout int64  `json:"shutdown_timeout_s"`
		ReadTimeout     int64  `json:"read_timeout_s"`
		WriteTimeout    int64  `json:"write_timeout_s"`
		IdleTimeout     int64  `json:"idle_timeout_s"`
	}

	// ConfigURL wrapper for net/url.URL
	ConfigURL url.URL

	// Service represents configuration of the service
	Service struct {
		VirtualHost string    `json:"virtual_host"`
		URL         ConfigURL `json:"url"`
		User        string    `json:"user"`
		Password    string    `json:"password"`
		RateLimit   int       `json:"rate_limit"`
		JWT         string    `json:"jwt"`
	}
)

// UnmarshalJSON implement unmarhsaling for URL
func (cu *ConfigURL) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	t := ConfigURL(*u)
	*cu = t

	return nil
}

// Load loads config  from file
func Load(file string) (Config, error) {
	var config Config

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
