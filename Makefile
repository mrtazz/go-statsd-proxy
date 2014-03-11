#
# simple makefile to run and build things
#
PROJECT=github.com/mrtazz/go-statsd-proxy

test:
	@echo "Running tests..."
	@go test ${PROJECT}/statsdproxy

benchmark:
	@echo "Running tests..."
	@go test -bench=. ${PROJECT}/statsdproxy

