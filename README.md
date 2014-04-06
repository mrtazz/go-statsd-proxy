# go-statsd-proxy

## Overview
A proxy for multiple statsd backends that routes metrics to specific instances
via consistent hashing. This is basically a reimplementation of the proxy
[included in Etsy's StatsD][statsd-proxy] and serves as a side project for me
to learn Go.

## Bugs
Probably a lot, submit them
[here](https://github.com/mrtazz/go-statsd-proxy/issues).

## Contributing
Take a look at [the
guidelines](https://github.com/mrtazz/go-statsd-proxy/blob/master/CONTRIBUTING.md).


[statsd-proxy]: https://github.com/etsy/statsd/blob/master/proxy.js
