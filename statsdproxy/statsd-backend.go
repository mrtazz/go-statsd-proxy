// definitions and functions for interacting with backends
package statsdproxy

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var healthCheckInterval int

type StatsDBackend struct {
	Host           string
	Port           int
	ManagementPort int
	conn           net.Conn
	ManagementConn net.Conn
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
	client.OpenManagementConnection()
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

// method to open TCP connection to the management port
func (client *StatsDBackend) OpenManagementConnection() {
	connectionString := fmt.Sprintf("%s:%d", client.Host, client.ManagementPort)
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		log.Println(err)
	}
	client.ManagementConn = conn
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
	if DebugMode {
		log.Printf("checking backend on %s:%d", client.Host, client.Port)
	}
	update_string := fmt.Sprintf("health")
	_, err := fmt.Fprintf(client.ManagementConn, update_string)
	if err != nil {
		log.Println(err)
	}
	reply := make([]byte, 1024)

	_, err = client.ManagementConn.Read(reply)
	if err != nil {
		log.Printf("Write to server failed:", err.Error())
	}
	health_status := strings.Trim(string(reply), string(0))

	if DebugMode {
		log.Printf("Response from backend %s:%d: %s", client.Host,
			client.ManagementPort, health_status)
	}
	if strings.Contains(health_status, "up") {
		if DebugMode {
			log.Printf("backend at %s:%d is up", client.Host, client.ManagementPort)
		}
		return true
	} else {
		if DebugMode {
			log.Printf("backend at %s:%d is down", client.Host, client.ManagementPort)
		}
		return false
	}
}

// function to figure out whether or not a backend is still up
//
// this queries the management interface of the backend to determine health
//
// returns false or true
func (client *StatsDBackend) Alive() bool {
	now := time.Now().Unix()
	if (now - client.Status.LastPingTime) > int64(healthCheckInterval) {
		if DebugMode {
			log.Printf("Checking alive status on backend %s:%d", client.Host,
				client.ManagementPort)
		}
		client.Status.Alive = client.CheckAliveStatus()
		client.Status.LastPingTime = now
	}
	return client.Status.Alive
}
