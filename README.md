# <img src="https://cloud.githubusercontent.com/assets/53900/26097013/7c930660-3a66-11e7-9b5c-780b0630d5a4.gif" alt="Muxy Logo" style="height: 80px;" height="80px"/>


Proxy for simulating real-world distributed system failures to improve resilience in your applications.

[![wercker status](https://app.wercker.com/status/e45703ebafd48632db56f022cc54546b/s "wercker status")](https://app.wercker.com/project/bykey/e45703ebafd48632db56f022cc54546b)
[![Go Report Card](https://goreportcard.com/badge/github.com/mefellows/muxy)](https://goreportcard.com/report/github.com/mefellows/muxy)
[![GoDoc](https://godoc.org/github.com/mefellows/muxy?status.svg)](https://godoc.org/github.com/mefellows/muxy)
[![Coverage Status](https://coveralls.io/repos/github/mefellows/muxy/badge.svg?branch=HEAD)](https://coveralls.io/github/mefellows/muxy?branch=HEAD)

## Introduction

Muxy is a proxy that _mucks_ with your system and application context, operating at Layers 4, 5 and 7, allowing you to simulate common failure scenarios from the perspective of an application under test; such as an API or a web application.

If you are building a distributed system, Muxy can help you test your resilience and fault tolerance patterns.

<p align="center">
  <img width="880" src="https://cdn.rawgit.com/mefellows/muxy/master/images/muxy.svg">
</p>

### Contents
<!-- TOC depthFrom:2 depthTo:4 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Introduction](#introduction)
	- [Contents](#contents)
- [Features](#features)
- [Installation](#installation)
	- [On Mac OSX using Homebrew](#on-mac-osx-using-homebrew)
	- [Using Go Get](#using-go-get)
- [Using Muxy](#using-muxy)
	- [5-minute example](#5-minute-example)
	- [Muxy as part of a test suite](#muxy-as-part-of-a-test-suite)
	- [Notes](#notes)
- [Proxies and Middlewares](#proxies-and-middlewares)
	- [Proxies](#proxies)
		- [HTTP Proxy](#http-proxy)
		- [TCP Proxy](#tcp-proxy)
	- [Middleware](#middleware)
		- [Delay](#delay)
		- [HTTP Tamperer](#http-tamperer)
		- [Network Shaper](#network-shaper)
		- [TCP Tamperer](#tcp-tamperer)
		- [Logger](#logger)
- [Configuration Reference](#configuration-reference)
- [Examples](#examples)
	- [Hystrix](#hystrix)
- [Usage with Docker](#usage-with-docker)
- [Extending Muxy](#extending-muxy)
	- [Proxies](#proxies)
	- [Middleware](#middleware)
- [Contributing](#contributing)

<!-- /TOC -->

## Features

* Ability to tamper with network devices at the transport level (Layer 4)
* Ability to tamper with the TCP session layer (Layer 5)
* ...and HTTP requests/responses at the HTTP protocol level (Layer 7)
    * Supports custom proxy routing (aka basic reverse proxy)
    * Advanced matching rules allow you to target specific requests
    * Introduce randomness into symptoms
* Simulate real-world network connectivity problems/partitions for mobile devices, distributed systems etc.
* Ideal for use in CI/Test Suites to test resilience across languages/technologies
* Simple native binary installation with no dependencies
* Extensible and modular architecture
* An official Docker [container](https://github.com/mefellows/docker-muxy) to simplify uses cases such as Docker Compose

## Installation

Download a [release](https://github.com/mefellows/muxy/releases) for your platform
and put it somewhere on the `PATH`.

### On Mac OSX using Homebrew

If you are using [Homebrew](http://brew.sh) you can follow these steps to install Muxy:

```bash
brew install https://raw.githubusercontent.com/mefellows/muxy/master/scripts/muxy.rb
```

### Using Go Get

```
go get github.com/mefellows/muxy
```

## Using Muxy

Muxy is typically used in two ways:

  1. In local development to see how your application responds
  under certain conditions
  1. In test suites to automate resilience testing

### 5-minute example

1. Install Muxy
1. Create configuration file `config.yml`:

    ```yaml
    # Configures a proxy to forward/mess with your requests
    # to/from www.onegeek.com.au. This example adds a 5s delay
    # to the response.
    proxy:
      - name: http_proxy
        config:
          host: 0.0.0.0
          port: 8181
          proxy_host: www.onegeek.com.au
          proxy_port: 80

    # Proxy plugins
    middleware:
      - name: http_tamperer
        config:
          request:
            host: "www.onegeek.com.au"

      # Message Delay request/response plugin
      - name: delay
        config:
          request_delay: 1000
          response_delay: 500

      # Log in/out messages
      - name: logger

    ```
1. Run Muxy with your config: `muxy proxy --config ./config.yml`
1. Make a request to www.onegeek.com via the proxy: `time curl -v -H"Host: www.onegeek.com.au" http://localhost:8181/`. Compare that with a request direct to the website: `time curl -v www.onegeek.com.au` - it should be approximately 5s faster.

That's it - running Muxy is a matter of configuring one or more [Proxies](#proxies), with 1 or more [Middleware](#middleware) components defined in a simple [YAML file](/examples/config.yml).

### Muxy as part of a test suite

1. Create an application
2. Build in fault tolerence (e.g. using something like [Hystrix](https://github.com/Netflix/Hystrix))
3. Create integration tests
  1. Run Muxy configuring a *proxy* such as HTTP, and one or more *symptom*s such as network latency, partition or HTTP error
  2. Point your app at Muxy
  3. Run tests and check if system behaved as expected
4. Profit!

### Notes

Muxy is a stateful system, and mucks with your low-level (system) networking interfaces and therefore cannot be run in parallel with other tests.
It is also recommended to run within a container/virtual machine to avoid unintended consequences (like breaking Internet access from the host).

## Proxies and Middlewares

### Proxies
#### HTTP Proxy

Simple HTTP(s) Proxy that starts up on a local IP/Hostname and Port.

Example configuration snippet:

```yaml
proxy:
  - name: http_proxy
    config:
      ## Proxy host details
      host: 0.0.0.0
      protocol: http
      port: 8181

      ## Proxy target details
      proxy_host: 0.0.0.0
      proxy_port: 8282
      proxy_protocol: https

      ## Certificate to present to Muxy clients (i.e. server certs)
      proxy_ssl_key: proxy-server/test.key
      proxy_ssl_cert: proxy-server/test.crt

      ## Certificate to present to Muxy proxy targets (i.e. client certs)
      proxy_client_ssl_key: client-certs/cert-key.pem
      proxy_client_ssl_cert: client-certs/cert.pem
      proxy_client_ssl_ca: client-certs/ca.pem

      ## Enable this to proxy targets we don't trust
      # insecure: true # allow insecure https

      # Specify additional proxy rules. Default catch-all proxy still
      # applies with lowest matching precedence.
      # Request matchers are specified as valid regular expressions
      # and must be properly YAML escaped.
      # See https://github.com/mefellows/muxy/issues/11 for behaviour.
      - request:
          method: 'GET|DELETE'
          path: '^\/foo'
          host: '.*foo\.com'
        pass:
          path: '/bar'
          scheme: 'http'
          host: 'bar.com'

```

#### TCP Proxy

Simple TCP Proxy that starts up on a local IP/Hostname and Port, forwarding traffic to the specified `proxy_host` on `proxy_port`.

Example configuration snippet:

```yaml
proxy:
  - name: tcp_proxy
    config:
      host: 0.0.0.0           # Local ip/hostname to bind to and accept connections.
      port: 8080              # Local port to bind to
      proxy_host: 0.0.0.0
      proxy_port: 2000
      nagles_algorithm: true
      packet_size: 64
```

### Middleware

Middleware have the ability to intervene upon receiving a request (Pre-Dispatch) or before sending the response back to the client (Post-Dispatch).
In some cases, such as the Network Shaper, the effect is applied _before any request is made_ (e.g. if the local network device configuration is altered).

#### Delay

A basic middleware that simply adds a delay of `delay` milliseconds to the request
or response.

Example configuration snippet:

```yaml
middleware:
  - name: delay
    config:
      request_delay: 1000      # Delay in ms to apply to request to target
      response_delay: 500      # Delay in ms to apply to response from target

      # Specify additional matching rules. Default is to apply delay to all
      # requests on all http proxies.
      # Request matchers are specified as valid regular expressions
      # and must be properly YAML escaped.
      # See https://github.com/mefellows/muxy/issues/11 for behaviour.
      matching_rules:
      - method: 'GET|DELETE'
        path: '^/boo'
        host: 'foo\.com'
```

#### HTTP Tamperer

A Layer 7 tamperer, this plugin allows you to modify response headers, status code or the body itself.

Example configuration snippet:

```yaml
middleware:
  - name: http_tamperer
    config:
      request:
        host:       "somehost"   # Override Host header that's sent to target
        path:             "/"    # Override the request path
        method:           "GET"  # Override request method
        headers:
          x_my_request:   "foo"  # Override request header
          content_type: "application/x-www-form-urlencoded"
          content_length: "5"
        cookies:                 # Custom request cookies
            - name: "fooreq"
              value: "blahaoeuaoeu"
              domain: "localhost"
              path: "/foopath"
              secure: true
              rawexpires: "Sat, 12 Sep 2015 09:19:48 UTC"
              maxage: 200
              httponly: true
        body: "wow, new body!"   # Override request body
      response:
        status: 201              # Override HTTP Status code
        headers:                 # Override response headers
          content_length: "27"
          x_foo_bar:      "baz"
        body:      "my new body" # Override response body
        cookies:                 # Custom response cookies
            - name: "foo"
              value: "blahaoeuaoeu"
              domain: "localhost"
              path: "/foopath"
              secure: true
              rawexpires: "Sat, 12 Sep 2015 09:19:48 UTC"
              maxage: 200
              httponly: true

      # Specify additional matching rules. Default is to apply delay to all
      # requests on all http proxies.
      # Request matchers are specified as valid regular expressions
      # and must be properly YAML escaped.
      # See https://github.com/mefellows/muxy/issues/11 for behaviour.
      matching_rules:
      - method: 'GET|DELETE'
        path: '^/boo'
        host: 'foo\.com'              
```

#### Network Shaper

The network shaper plugin is a Layer 4 tamperer, and requires *root access* to work, as it needs to configure the local firewall and network devices.
Using the excellent [Comcast](https://github.com/tylertreat/comcast) library, it can shape and interfere with network traffic,
including bandwidth, latency, packet loss and jitter on specified ports, IPs and protocols.

NOTE: This component only works on MacOSX, FreeBSD, Linux and common *nix flavours.

Example configuration snippet:

```yaml
middleware:

  - name: network_shape
    config:
      latency:     250         # Latency to add in ms
      target_bw:   750         # Bandwidth in kbits/s
      packet_loss: 0.5         # Packet loss, as a %
      target_ips:              # Target ipv4 IP addresses
        - 0.0.0.0
      target_ips6:             # Target ipv6 IP addresses
        - "::1/128"
      target_ports:            # Target destination ports
        - "80"
      target_protos:           # Target protocols
        - "tcp"
        - "udp"
        - "icmp"
```

#### TCP Tamperer

The TCP Tamperer is a Layer 5 tamperer, modifying the messages in and around TCP
sessions. Crudely, you can set the body of inbound and outbound TCP packets, truncate
the last character of messages or randomise the text over the wire.

```yaml
- name: tcp_tamperer
  config:
    request:
      body: "wow, new request!"   # Override request body
      randomize: true             # Replaces input message with a random string
      truncate: true              # Removes last character from the request message
    response:
      body: "wow, new response!" # Override response body
      randomize: true             # Replaces response message with a random string
      truncate: true              # Removes last character from the response message
```

#### Logger

Log the in/out messages, optionally requesting the output to be hex encoded.

Example configuration snippet:

```yaml
middleware:
  - name: logger
    config:
      hex_output: false        # Display output as Hex instead of a string
```

## Configuration Reference

Refer to the [example](/examples/config.yml) YAML file for a full reference.

## Examples

### Hystrix

Using the [Hystrix Go](https://github.com/afex/hystrix-go) library, we use Muxy to trigger a circuit breaker and return a canned response, ensuring we don't have downtime. View the [example](examples/hystrix).

## Usage with Docker

Download the [Docker image](https://github.com/mefellows/docker-muxy) by running:

```
docker pull mefellows/muxy
```

After creating a [config](#configuration-reference] file (let's assume it's at `./conf/config.yml`), and assuming you are proxying something on port `80`, you can now run the image locally:

```
docker run \
  -d \
  -p 80:80 \
  -v "$PWD/conf":/opt/muxy/conf \
  --privileged \
  mefellows/muxy
```

You should now be able to hit this Docker container and simulate any failures as per usual. e.g. `curl docker:80/some/endpoint`.

The [Hystrix](#hystrix) example above has a detailed example on how to use Muxy with a more complicated system, using Docker Compose to orchestrate a number of containers.

## Extending Muxy

Muxy is built as a series of configurable plugins (using [Plugo](https://github.com/mefellows/plugo)) that must be specified in the configuration
file to be activated. Start with a quick tour of Plugo before progressing.

### Proxies

Proxies must implement the [Proxy](/muxy/proxy.go) interface, and register themselves via `PluginFactories.register` to be available at runtime.
Take a look at the [HTTP Proxy](protocol/http.go) for a good working example.

### Middleware

Middlewares implement the [Middleware](/muxy/middle.go) interface  and register themselves via `PluginFactories.register` to be available at runtime.
Take a look at the [HTTP Delay](symptom/http_delay.go) for a good working example.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).
