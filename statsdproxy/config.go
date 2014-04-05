// parser for the proxy config json file

package statsdproxy

import (
	"encoding/json"
	"io/ioutil"
)

// struct to represent a backend node in JSON
type StatsdConfigNode struct {
	Host      string `json:host`
	Port      int    `json:port`
	Adminport int    `json:adminport`
}

// struct to represent the whole JSON config file
type ProxyConfig struct {
	Host          string             `json:host`
	Port          int                `json:port`
	CheckInterval int                `json:checkInterval`
	Nodes         []StatsdConfigNode `json:nodes`
}

// constructor function to create a new config struct with values
// accepts a filepath to a config file as parameter
// returns the config struct and an error
func NewConfig(filepath string) (*ProxyConfig, error) {
	raw_config, err := ioutil.ReadFile(filepath)
	config, err := readConfigFile(raw_config)
	return config, err
}

// function to parse the raw json data into a ProxyConfig struct
// accepts a raw byte array as parameter
// returns the ProxyConfig struct and error
func readConfigFile(data []byte) (*ProxyConfig, error) {
	var config ProxyConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
