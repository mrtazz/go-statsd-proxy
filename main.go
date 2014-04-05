// main module to startup the statsd proxy
package main

import (
	"flag"
	"github.com/mrtazz/go-statsd-proxy/statsdproxy"
)

func main() {
	// note, that variables are pointers
	configfile := flag.String("f", "", "config file path")
	flag.Parse()
	statsdproxy.StartListener(*configfile)
}
