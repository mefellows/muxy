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