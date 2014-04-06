# go-statsd-proxy

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
