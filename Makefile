.PHONY: build clean

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/passenger-datadog-monitor *.go

clean:
	rm -rf ./bin