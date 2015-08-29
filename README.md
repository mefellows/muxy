# Muxy

Simulating real-world distributed system failures to improve resilience in your applications.

[![wercker status](https://app.wercker.com/status/e45703ebafd48632db56f022cc54546b/s "wercker status")](https://app.wercker.com/project/bykey/e45703ebafd48632db56f022cc54546b)

## Introduction

Muxy is a proxy that _mucks_ around with your system and application context, allowing you to simulate common failure scenarios from the perspective of an application under test - such as an API or a web application. 

If you are building a distributed system, Muxy can help you test your resilience and fault tolerance patterns

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
