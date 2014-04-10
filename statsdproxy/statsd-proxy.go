// runner to start up the proxy
package statsdproxy

import (
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	CHANNEL_SIZE = 100
)

// variable to indicate whether or not we run in DebugMode
var DebugMode bool

// channel to gather internal metrics
var internalMetrics chan StatsDMetric
var metricsOutput chan metricsRequest

type StatsDMetric struct {
	name  string
	value float64
	raw   string
}

// exported functions
func StartProxy(cfgFilePath string, quit chan bool) error {
	var err error
	config, err := NewConfig(cfgFilePath)
	if err != nil {
		log.Printf("Error parsing config file: %s (exiting...)", err)
		return nil
	}
	internalMetrics = make(chan StatsDMetric, CHANNEL_SIZE)
	metricsOutput = make(chan metricsRequest, CHANNEL_SIZE)

	go metricsCollector(internalMetrics)
	go StartMainListener(*config)
	go StartManagementConsole(*config)

	// wait until you're told to quit
	<-quit

	return nil

}

// function to set up the main UDP listener. Everything that is needed to
// receive and relay metrics
func StartMainListener(config ProxyConfig) error {

	log.Printf("Starting StatsD listener on %s and port %d", config.Host, config.Port)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(config.Host), Port: config.Port})
	if err != nil {
		log.Printf("Error setting up listener: %s (exiting...)", err)
		return nil
	}

	relay_channel := make(chan StatsDMetric, CHANNEL_SIZE)
	hash_ring := *NewHashRing()
	for _, node := range config.Nodes {
		backend := NewStatsDBackend(node.Host, node.Port, node.Adminport,
			config.CheckInterval)
		if DebugMode {
			log.Printf("Adding backend %s:%d", backend.Host, backend.Port)
		}
		hash_ring, err = hash_ring.Add(*backend)
		if err != nil {
			log.Println("Error adding backend to Hashring")
			log.Println(err)
		}
	}
	go relay_metric(hash_ring, relay_channel)

	for {
		buf := make([]byte, 512)
		num, _, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP buffer: %s (skipping...)", err)
			return nil
		} else {
			go handleConnection(buf[0:num], relay_channel)
		}
	}

	return nil
}

// handle the actual incoming connection and figure out which packet types we
// got sent.
// accepts a byte array of data
func handleConnection(data []byte, relay_channel chan StatsDMetric) {
	if DebugMode {
		log.Printf("Got packet: %s", string(data))
	}
	metrics := strings.Split(string(data), "\n")
	for _, str := range metrics {
		metric := parsePacketString(str)
		internalMetrics <- *metric
		relay_channel <- *metric
	}

}

// parse a string into a statsd packet
// accepts a string of data
// returns a StatsDMetric
func parsePacketString(data string) *StatsDMetric {
	ret := new(StatsDMetric)
	first := strings.Split(data, ":")
	if len(first) < 2 {
		log.Printf("Malformatted metric: %s", data)
		return ret
	}
	name := first[0]
	second := strings.Split(first[1], "|")
	value64, _ := strconv.ParseInt(second[0], 10, 0)
	value := float64(value64)
	// check for a samplerate
	third := strings.Split(second[1], "@")
	metric_type := third[0]

	switch metric_type {
	case "c", "ms", "g":
		ret.name = name
		ret.value = value
		ret.raw = data
	default:
		log.Printf("Unknown metrics type: %s", metric_type)
	}

	return ret
}

// relay a metric to one of the active statsd backends
func relay_metric(ring HashRing, relay_channel chan StatsDMetric) {
	for {
		select {
		case metric := <-relay_channel:
			// find out which backend to relay to and do it
			backend_host, err := ring.GetBackendForMetric(metric.name)
			if err != nil {
				log.Printf("Unable to get backend for metric: %s", metric.name)
			} else {
				if DebugMode {
					log.Printf("relaying metric: %s to %s:%d", metric.raw,
						backend_host.Host, backend_host.Port)
				}
				backend_host.Send(metric.raw)
			}
		}
	}
}
