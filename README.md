# Muxy

Simulating real-world distributed system failures to improve resilience in your applications.

[![wercker status](https://app.wercker.com/status/e45703ebafd48632db56f022cc54546b/s "wercker status")](https://app.wercker.com/project/bykey/e45703ebafd48632db56f022cc54546b)

## Introduction

Muxy is a proxy that _mucks_ with your system and application context, operating at Layers 4 and 7, allowing you to simulate common failure scenarios from the perspective of an application under test; such as an API or a web application.

If you are building a distributed system, Muxy can help you test your resilience and fault tolerance patterns.

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

That's it!

### Muxy as part of a test suite

T

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
