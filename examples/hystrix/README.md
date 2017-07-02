# Hystrix Example

A simple Docker Compose setup that has a frontend API <- issues requests to -> backend API. In the event that:

1. The API returns an error (40x, 50x response code)
1. The API is inaccessible
1. The API does not meet its defined SLAs (currently a generous 1500ms)

The API should open a [circuit breaker](https://github.com/Netflix/Hystrix/wiki/How-it-Works#CircuitBreaker) and return a canned response, instead of querying the back-end. We use the [Go Hystrix](https://github.com/afex/hystrix-go) library to achieve this, along with Muxy to interfere and trigger this behaviour.

## Running the example

Ensure that Docker and Docker Compose is installed, and then run:

```
./run-tests.sh
```

## Run manually

Each of the following in separate tabs:

```
cd api && API_HOST=http://localhost:8001 PORT=8000 STATSD_HOST=192.168.99.100:8125 go run main.go
muxy proxy --config muxy/conf/config.local.yml
cd backup && PORT=8002 STATSD_HOST=192.168.99.100:8125 go run main.go
time echo "GET http://localhost:8000/" | vegeta attack -duration=15s | tee results.bin | vegeta report
```

## Hystrix Dashboard

### Local API
Grab your own IP address (e.g. 192.168.0.7)

http://docker:7979/hystrix-dashboard/monitor/monitor.html?streams=%5B%7B%22name%22%3A%22Test%22%2C%22stream%22%3A%22http%3A%2F%2F192.168.0.7%3A8181%22%2C%22auth%22%3A%22%22%2C%22delay%22%3A%22%22%7D%5D

### Docker API

