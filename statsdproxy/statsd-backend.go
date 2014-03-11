// definitions and functions for interacting with backends
package statsdproxy

import (
	"log"
	"net"
	"fmt"
  "time"
)

const (
  PING_CHECK_INTERVAL = 10
)

type StatsDBackend struct {
	Host string
	Port int
	ManagementPort int
	conn net.Conn
	RingID HashRingID
	Status struct {
	  Alive bool
	  LastPingTime int64
  }
}

func NewStatsDBackend(host string, port int) *StatsDBackend {
	client := StatsDBackend{Host: host, Port: port}
  client.RingID, _ = GetHashRingPosition(fmt.Sprintf("%s:%s", host, port))
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
  if ( now - client.Status.LastPingTime) > PING_CHECK_INTERVAL {
    client.Status.Alive = client.CheckAliveStatus()
    client.Status.LastPingTime = now
  }
  return client.Status.Alive
}
