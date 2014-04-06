// this file contains functions and helpers for collecting overall metrics and
// running a management interface

package statsdproxy

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	RECV_BUF_LEN = 1024
)

// start the management console on all interface on the configured port
func StartManagementConsole(config ProxyConfig) error {

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d",
		config.ManagementPort))
	if err != nil {
		log.Printf("error listening: %s", err.Error())
		return err
	}

	buf := make([]byte, RECV_BUF_LEN)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accept: %s", err.Error())
			return err
		}
		_, err = conn.Read(buf)
		if err != nil {
			log.Printf("Error reading:", err.Error())
			return err
		}
		answer := answerManagementQuery(strings.Trim(string(buf), string(0)))
		_, err = conn.Write([]byte(answer))
		if err != nil {
			log.Printf("Error writing to connection: %s", err.Error())
		}
	}
	return nil
}

// all command answering is done in this function
func answerManagementQuery(query string) string {
	var answer string
	query = strings.Trim(query, " \n\r")
	switch query {
	case "ping":
		answer = "pong"
	default:
		if DebugMode {
			log.Printf("nothing known for command %s", query)
		}
		answer = "unknown command"
	}
	return answer
}
