// definitions and functions for interacting with backends
package statsdproxy

import (
	"fmt"
	"log"
	"net"
	"time"
)

var healthCheckInterval int

type StatsDBackend struct {
	Host           string
	Port           int
	ManagementPort int
	conn           net.Conn
	RingID         HashRingID
	Status         struct {
		Alive        bool
		LastPingTime int64
	}
}

func NewStatsDBackend(host string, port int,
	managementPort int, check_interval int) *StatsDBackend {
	healthCheckInterval = check_interval
	client := StatsDBackend{Host: host, Port: port, ManagementPort: managementPort}
	client.RingID, _ = GetHashRingPosition(fmt.Sprintf("%s:%d", host, port))
	client.Open()
	return &client
}

// Method to open udp connection, called by default client factory
func (client *StatsDBackend) Open() {
	connectionString := fmt.Sprintf("%s:%d", client.Host, client.Port)
	conn, err := net.Dial("udp", connectionString)
	if err != nil {
		log.Println(err)
	}
	client.conn = conn
}

// Method to close udp connection
func (client *StatsDBackend) Close() {
	client.conn.Close()
}

func (client *StatsDBackend) Send(data string) {
	if DebugMode {
		log.Printf("sending %s to backend on port %d", data, client.Port)
	}
	update_string := fmt.Sprintf(data)
	_, err := fmt.Fprintf(client.conn, update_string)
	if err != nil {
		log.Println(err)
	}
}

func (client *StatsDBackend) CheckAliveStatus() bool {
	// TODO: check management console of client
	return true
}

// function to figure out whether or not a backend is still up
//
// this queries the management interface of the backend to determine health
//
// returns false or true
func (client *StatsDBackend) Alive() bool {
	now := time.Now().Unix()
	if (now - client.Status.LastPingTime) > int64(healthCheckInterval) {
		client.Status.Alive = client.CheckAliveStatus()
		client.Status.LastPingTime = now
	}
	return client.Status.Alive
}
