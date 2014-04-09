// main module to startup the statsd proxy
package main

import (
	"flag"
	"github.com/mrtazz/go-statsd-proxy/statsdproxy"
)

func main() {
	// note, that variables are pointers
	configfile := flag.String("f", "", "config file path")
	debug := flag.Bool("d", false, "enable debug mode")
	flag.Parse()
	statsdproxy.DebugMode = *debug
  // TODO: catch sigterm and actually tell the process to quit
  quit := make(chan bool)
	statsdproxy.StartProxy(*configfile, quit)
}
