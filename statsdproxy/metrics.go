// this file contains functions and helpers for collecting overall metrics and
// running a management interface

package statsdproxy

import (
	"fmt"
	"log"
	"net"
	"strings"
	"runtime"
)

const (
	RECV_BUF_LEN = 1024
)

type Answers []string
type MetricsCollection map[string]float64
type metricsRequest struct {
  response chan []string
}

// start the management console on all interface on the configured port
func StartManagementConsole(config ProxyConfig) error {
	log.Printf("Starting Management listener on %s and port %d", config.Host,
	config.ManagementPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d",
		config.ManagementPort))
	if err != nil {
		log.Printf("error listening: %s", err.Error())
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accept: %s", err.Error())
			return err
		}
    buf := make([]byte, RECV_BUF_LEN)
		_, err = conn.Read(buf)
		if err != nil {
			log.Printf("Error reading:", err.Error())
			return err
		}
		answers := answerManagementQuery(strings.Trim(string(buf), string(0)))

    for _, value := range answers {
      _, err = conn.Write([]byte(fmt.Sprintf("%s\n", value)))
      if err != nil {
        log.Printf("Error writing to connection: %s", err.Error())
      }
    }
	}
	return nil
}

// all command answering is done in this function. We basically match commands
// to functions
func answerManagementQuery(query string) Answers {
	answer := make(Answers, 0, 20)
	query = strings.Trim(query, " \n\r")
	switch query {
	case "ping":
    answer = append(answer, "pong")
  case "memstats":
    answer = getMemStats()
  case "stats":
    result := make(chan []string)
    metricsOutput <- metricsRequest{make(chan []string)}
    answer = <- result
	default:
		if DebugMode {
			log.Printf("nothing known for command %s", query)
		}
    answer = append(answer, "unknown command")
	}
	return answer
}

// assemble memory stats in an array to print them out
func getMemStats() Answers {
	answer := make(Answers, 0, 20)
  memStats := new(runtime.MemStats)
  runtime.ReadMemStats(memStats)

  answer = append(answer, fmt.Sprintf("bytes allocated and in use: %d", memStats.Alloc))
  answer = append(answer, fmt.Sprintf("bytes allocated total: %d", memStats.TotalAlloc))
  answer = append(answer, fmt.Sprintf("bytes obtained from system: %d", memStats.Sys))
  answer = append(answer, fmt.Sprintf("number of pointer lookups: %d", memStats.Lookups))
  answer = append(answer, fmt.Sprintf("number of mallocs: %d", memStats.Mallocs))
  answer = append(answer, fmt.Sprintf("number of frees: %d", memStats.Frees))
  answer = append(answer, fmt.Sprintf("bytes allocated and still in use: %d", memStats.HeapAlloc))
  answer = append(answer, fmt.Sprintf("bytes obtained from system: %d", memStats.HeapAlloc))
  answer = append(answer, fmt.Sprintf("bytes in idle spans: %d", memStats.HeapIdle))
  answer = append(answer, fmt.Sprintf("bytes in non-idle span: %d", memStats.HeapInuse))
  answer = append(answer, fmt.Sprintf("bytes released to the OS: %d", memStats.HeapReleased))
  answer = append(answer, fmt.Sprintf("total number of allocated objects: %d", memStats.HeapObjects))

  return answer
}

// gather internal metrics
func metricsCollector(metricsChannel chan StatsDMetric) {

  metrics := make(MetricsCollection)

	for {
		select {
		case metric := <- metricsChannel:
		  if DebugMode {
        log.Printf("Received internal metric: %s", metric.raw)
      }
      if _,ok := metrics[metric.name]; ok {
        metrics[metric.name] += metric.value
      } else {
        metrics[metric.name] = metric.value
      }
    case request := <- metricsOutput:
      var out []string
      for key, value := range metrics {
        out = append(out, fmt.Sprintf("%s: %f", key, value))
      }
      request.response <- out
		}
	}
}
