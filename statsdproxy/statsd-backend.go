// definitions and functions for interacting with backends
package statsdproxy

import (
	"log"
	"net"
	"fmt"
)

type StatsDBackend struct {
	Host string
	Port int
	conn net.Conn
	RingID HashRingID
}

func New(host string, port int) *StatsDBackend {
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
