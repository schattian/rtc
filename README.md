# RTC

[![Build Status](https://travis-ci.com/sebach1/rtc.svg?branch=master)](https://travis-ci.com/sebach1/rtc)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebach1/rtc)](https://goreportcard.com/report/github.com/sebach1/rtc)

**DISCLAIMER:** The project is currently in WIP state. Maybe isn't yet good time to contribute due things will keep moving around all the time.

## Overview

RTC is a framework to call remote transactions over multiple services.

It focuses on be:

- **Centralized:** transactions will be executed following a unique pipeline.

- **Cross-language:** providing APIs over HTTP / gRPC to use it and also config via YAML files.

- **No-protocol limited:** protocols of services are not a limitation at the hour of integrating any of them. For more, see the [list of supported protocols].

- **Simple:** to ensure simplicity over differences across services and provide a normalized way to interact with data pipelines, it uses [git abstractions]() and a [schema]().

- **No ownership needed:** the services to be used are all treated as foreign ones, so there isn't a client-side on any RTC callable service. You call it by its own interface.

## Purpose

The original and minimal purpose is to forget SDKs.
That means, removing:

- Huge amount of dependencies when connecting multiple services.

- Maintainability of the SDK itself.

- Reliability on ~~never~~ not always tested software.

- Lock of the mainstream-only technology stack.

The second one is to improve the understanding and knowledge of external data structures.

It comes from the Rob Pike's rule:

> *Data dominates. If you've chosen the right data structures and organized things well, the algorithms will almost always be self-evident. Data structures, not algorithms, are central to programming.*

Its very harmful (and then, a frequent mistake) to choose bad data structures by misunderstanding the external structures in which your act will depend on.

## Usage

In case you use Go, its recommendable to inspect the [godoc] for a complete usage guide omitting the usage of an API.
