# go-statsd-proxy [![Build Status](https://travis-ci.org/mrtazz/go-statsd-proxy.svg?branch=master)](https://travis-ci.org/mrtazz/go-statsd-proxy)

## Overview
A proxy for multiple statsd backends that routes metrics to specific instances
via consistent hashing. This is basically a reimplementation of the proxy
[included in Etsy's StatsD][statsd-proxy] and serves as a side project for me
to learn Go.

## Usage
```
git clone https://github.com/mrtazz/go-statsd-proxy
cd go-statsd-proxy
go build
./go-statsd-proxy -f exampleConfig.json
```

## Monitoring
The proxy has a management interface accessible via TCP (inspired by the
StatsD interface) which can be used for monitoring and accessing some stats
about the running process. By default the interface runs on port 8126.

### Ping
This can be used as a basic health check to see if the server is still
responding. It's not really detailed or granular but may change in the future.
```
% echo "ping" | nc -w1 localhost 8126                                                                                                                                                           <master ✗>
pong
```

### Stats
This command gives you an overview over some of the internal stats of the
running proxy:

```
% echo -n "stats" | nc -w1 127.0.01 8126
time running in seconds: 51
packets_received: 1.000000
```

### Memstats
This command fives you an overview over the most important memory stats. Use
this to feed instance metrics into ganglia for example:

```
% echo "memstats" | nc -w1 localhost 8126
bytes allocated and in use: 292432
bytes allocated total: 363088
bytes obtained from system: 4331752
number of pointer lookups: 1091
number of mallocs: 874
number of frees: 565
bytes allocated and still in use: 292432
bytes obtained from system: 292432
bytes in idle spans: 610304
bytes in non-idle span: 438272
bytes released to the OS: 610304
total number of allocated objects: 309
```

## Bugs
Probably a lot, submit them
[here](https://github.com/mrtazz/go-statsd-proxy/issues).

There is also a debug mode included which probably makes a lot of noise
depending on how many metrics you send. So be warned. It can be enabled by
running the proxy with the `-d` flag.

## Contributing
Take a look at [the
guidelines](https://github.com/mrtazz/go-statsd-proxy/blob/master/CONTRIBUTING.md).


[statsd-proxy]: https://github.com/etsy/statsd/blob/master/proxy.js
