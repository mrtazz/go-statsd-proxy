// main module to startup the statsd proxy
package main

import (
	"flag"
	"github.com/mrtazz/go-statsd-proxy/statsdproxy"
)

func main() {
	// note, that variables are pointers
	address := flag.String("a", "0.0.0.0", "address to listen on (default: 0.0.0.0")
	port := flag.Int("p", 8125, "port to listen on (default: 8125")
	statsdproxy.StartListener(*address, *port)
}
