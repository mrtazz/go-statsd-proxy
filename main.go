// main module to startup the statsd proxy
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrtazz/go-statsd-proxy/statsdproxy"
)

var (
	configfile = flag.String("f", "/etc/statsdproxy.json", "Configuration file path")
	debug      = flag.Bool("d", false, "enable debug mode")
)

func main() {
	flag.Parse()

	statsdproxy.DebugMode = *debug

	quit := make(chan bool)

	handleSignals(quit)

	statsdproxy.StartProxy(
		*configfile,
		quit,
	)
}

func handleSignals(quitChan chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func(chan os.Signal) {
		for sig := range c {
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				quitChan <- true
				os.Exit(0)
			}
		}
	}(c)
}
