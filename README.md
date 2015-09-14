# Muxy

Proxy for simulating real-world distributed system failures to improve resilience in your applications.

[![wercker status](https://app.wercker.com/status/e45703ebafd48632db56f022cc54546b/s "wercker status")](https://app.wercker.com/project/bykey/e45703ebafd48632db56f022cc54546b)

## Introduction

Muxy is a proxy that _mucks_ with your system and application context, operating at Layers 4 and 7, allowing you to simulate common failure scenarios from the perspective of an application under test; such as an API or a web application.

If you are building a distributed system, Muxy can help you test your resilience and fault tolerance patterns.

### Contents

  * [Features](#features)
  * [Installation](#installation)
  * [Using Muxy](#using-muxy)
    * [5 Minute Quick Start](#5-minute-example)
  * [Muxy Components](#proxies-and-middlewares)
    * [Proxies](#proxies)
      * [HTTP Proxy](#http-proxy)
      * [TCP Proxy](#tcp-proxy)
    * [Middleware](#middleware)
      * [HTTP Delay](#http-delay)
      * [HTTP Tamperer](#http-tamperer)
      * [Network Shaper](#network-shaper)
      * [Logger](#logger)
  * [YAML Configuration Reference](#configuration-reference)
  * [Extending Muxy](#extending-muxy)

## Features

* Ability to tamper with network devices at the transport level (Layer 4) 
* ...and HTTP requests/responses at the HTTP protocol level (Layer 7)
* Simulate real-world network connectivity problems/partitions for mobile devices, distributed systems etc.
* Ideal for use in CI/Test Suites to test resilience across languages/technologies
* Simple native binary installation with no dependencies
* Extensible and modular architecture

## Installation

Download a [release](https://github.com/mefellows/muxy/releases) for your platform
and put it somewhere on the `PATH`.

### On Mac OSX using Homebrew

If you are using [Homebrew](http://brew.sh) you can follow these steps to install Muck:

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
          proxy_host: onegeek.com.au
          proxy_port: 80

    # Proxy plugins
    middleware:

      # HTTP response delay plugin
      - name: http_delay
        config:
          delay: 5

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

Simple HTTP Proxy that starts up on a local IP/Hostname and Port. 

Example configuration snippet:

```yaml
proxy:
  - name: http_proxy
    config:
      host: 0.0.0.0
      protocol: http
      port: 8181
      proxy_host: 0.0.0.0
      proxy_port: 8282
      proxy_protocol: https
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

#### HTTP Delay

A basic middleware that simply adds a delay of `delay` seconds.

Example configuration snippet:

```yaml
middleware:
  - name: http_delay
    config:
      delay: 1                 # Delay in seconds to apply to response
```

#### HTTP Tamperer

A Layer 7 tamperer, this plugin allows you to modify response headers, status code or the body itself.

Example configuration snippet:

```yaml
middleware:
  - name: http_tamperer
    config:
      request:
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

## Extending Muxy

Muxy is built as a series of configurable plugins (using [Plugo](https://github.com/mefellows/plugo)) that must be specified in the configuration
file to be activated. Start with a quick tour of Plugo before progressing.

### Proxies

Proxies must implement the [Proxy](/muxy/proxy.go) interface, and register themselves via `PluginFactories.register` to be available at runtime. 
Take a look at the [HTTP Proxy](protocol/http.go) for a good working example.

### Middleware

Middlewares implement the [Middleware](/muxy/middle.go) interface  and register themselves via `PluginFactories.register` to be available at runtime. 
Take a look at the [HTTP Delay](symptom/http_delay.go) for a good working example.
