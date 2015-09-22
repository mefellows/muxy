# Golang Docker Echo Server

A small echo server that will respond to headers for the purposes of testing
the Nginx configuration.

## Endpoints

* `GET /headers/:header` - Echoes the requested `header` it receives back to you
* `GET /*` - Responds to any other request with `pong`

## Getting Started

* Install [Docker](http://docker.io/) and Docker Machine
* Install the Go cross-compiler [Gox](https://github.com/mitchellh/gox)

	```
	go get github.com/mitchellh/gox
	gox -build-toolchain
	```

## Build and Run application in Docker

```
gox -osarch="linux/amd64" -output="echoserver" && docker build -t mefellows/echoserver . && docker run --name api --rm -i -t -p 8000:8000 mefellows/echoserver
```

This will build a binary for 64bit linux with the name `echoserver`, create a Docker image and then run that image exposing port `8000`.

The API is now running on port 8000, `curl` until your heart is content:

```
curl $(docker-machine ip dev):8000/ping
```
