package statsdproxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func getJSONFixtureData(name string) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	raw_json, err := ioutil.ReadFile(pwd + "/test_fixtures/" + name + ".json")
	if err != nil {
		return nil, err
	}
	return raw_json, nil

}

func TestParsePacketString_withCorrectCounter(t *testing.T) {

	type Data struct {
		Raw        string
		Name       string
		Value      float64
		Samplerate float32
	}

	raw_json, err := getJSONFixtureData("correct_counters")
	if err != nil {
		t.Errorf("Error loading JSON data: %v", err)
	}

	var raw_data []Data
	if err := json.Unmarshal(raw_json, &raw_data); err != nil {
		t.Errorf("Error parsing JSON data: %v", err)
	}

	for _, entry := range raw_data {
		data := entry.Raw
		metric := parsePacketString(data)

		if metric.name != entry.Name {
			t.Errorf("parsePacketString: expected name %s and got %s", entry.Name, metric.name)
		}
		if metric.value != entry.Value {
			t.Errorf("parsePacketString: expected value %v and got %v", entry.Value, metric.value)
		}
	}
}

func TestParsePacketString_withCorrectTiming(t *testing.T) {
	type Data struct {
		Raw        string
		Name       string
		Value      float64
		Samplerate float32
	}

	raw_json, err := getJSONFixtureData("correct_timers")
	if err != nil {
		t.Errorf("Error loading JSON data: %v", err)
	}

	var raw_data []Data
	if err := json.Unmarshal(raw_json, &raw_data); err != nil {
		t.Errorf("Error parsing JSON data: %v", err)
	}

	for _, entry := range raw_data {
		data := entry.Raw
		metric := parsePacketString(data)

		if metric.name != entry.Name {
			t.Errorf("parsePacketString: expected name %s and got %s", entry.Name, metric.name)
		}
		if metric.value != entry.Value {
			t.Errorf("parsePacketString: expected value %v and got %v", entry.Value, metric.value)
		}
	}
}

func TestParsePacketString_withCorrectGauge(t *testing.T) {
	type Data struct {
		Raw        string
		Name       string
		Value      float64
		Samplerate float32
	}

	raw_json, err := getJSONFixtureData("correct_gauges")
	if err != nil {
		t.Errorf("Error loading JSON data: %v", err)
	}

	var raw_data []Data
	if err := json.Unmarshal(raw_json, &raw_data); err != nil {
		t.Errorf("Error parsing JSON data: %v", err)
	}

	for _, entry := range raw_data {
		data := entry.Raw
		metric := parsePacketString(data)

		if metric.name != entry.Name {
			t.Errorf("parsePacketString: expected name %s and got %s", entry.Name, metric.name)
		}
		if metric.value != entry.Value {
			t.Errorf("parsePacketString: expected value %v and got %v", entry.Value, metric.value)
		}
	}
}

// Benchmarks
func BenchmarkParsePacketString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePacketString("foo:1|c")
	}
}
