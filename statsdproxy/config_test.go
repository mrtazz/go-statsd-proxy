package statsdproxy

import (
	"testing"
)

func TestReadConfigFile_BasicData(t *testing.T) {
	const testConfig = `
    {
      "nodes": [
        {"host": "127.0.0.1", "port": 8129, "adminport": 8126},
        {"host": "127.0.0.1", "port": 8127, "adminport": 8128},
        {"host": "127.0.0.1", "port": 8129, "adminport": 8130}
      ],
      "host":  "0.0.0.0",
      "port": 8125,
      "checkInterval": 1000,
      "cacheSize": 10000
    }
    `

	parsed_config, err := readConfigFile([]byte(testConfig))

	if err != nil {
		t.Errorf("readConfigFile() parsing is broken with %v", err)
		t.FailNow()
	}
	if parsed_config.Port != 8125 {
		t.Errorf("wrong port read, expected 8125 and got %d", parsed_config.Port)
	}
	if parsed_config.Host != "0.0.0.0" {
		t.Errorf("wrong host read, expected 0.0.0.0 and got %s", parsed_config.Host)
	}
	if parsed_config.CheckInterval != 1000 {
		t.Errorf("wrong CheckInterval read, expected 1000 and got %d",
			parsed_config.CheckInterval)
	}
	if len(parsed_config.Nodes) != 3 {
		t.Errorf("wrong number of nodes read, expected 3 and got %d",
			len(parsed_config.Nodes))
	}
	node := parsed_config.Nodes[0]
	if node.Host != "127.0.0.1" {
		t.Errorf("wrong node Host read, expected 127.0.0.1 and got %s", node.Host)
	}
	if node.Port != 8129 {
		t.Errorf("wrong node Port read, expected 8129 and got %d", node.Port)
	}
	if node.Adminport != 8126 {
		t.Errorf("wrong node Adminport read, expected 8126 and got %d",
			node.Adminport)
	}

}
