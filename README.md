# RTC

[![Build Status](https://travis-ci.com/sebach1/rtc.svg?branch=master)](https://travis-ci.com/sebach1/rtc)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebach1/rtc)](https://goreportcard.com/report/github.com/sebach1/rtc)

  - [Overview](#overview)
  - [Purpose](#purpose)
  - [Usage](#usage)
  - [Concepts](#concepts)
    - [Schema](#schema)
    - [Git](#git)

**DISCLAIMER:** The project is currently in WIP state. Maybe isn't yet good time to contribute due things will keep moving around all the time.

## Overview

RTC is a framework to call remote transactions over multiple services.

It focuses on be:

- **Centralized:** transactions will be executed following a unique pipeline.

- **Cross-language:** providing APIs over HTTP / gRPC to use it and also config via YAML files.

- **No-protocol limited:** protocols of services are not a limitation at the hour of integrating any of them. For more, see the [list of supported protocols]().

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

In case you use Go, its recommendable to inspect the [godoc]() for a complete usage guide omitting the usage of an API.

## Concepts

### Schema

A schema describes the data structure of a callable transactional service.
The definition of it can be served across HTTPs (k8s like) or over a filesystem via YAML. For more, see the [Schema reference]().

It also has its own structure. Formed by a tiny sql abstraction:

- **Blueprint:** the group of tables a schema contains. It just describes all the structure, but it doesn't relate it.

- **Tables:** *sql-like* concept. It describes a collection of abstractions that are described by the same columns.

- **Columns:** *sql-like* concept. It holds a specific type of data (*type-safe*) and describes a field over an entity.
  
  - **Type:** limits the type a column can have.

- **Option:** just a key-value storage used to pass extra information about the transaction that can't be expressed through columns of an entity (e.g: scopes).

  - **Option key:** limits the possible key of the storage.

### Git

RTC provides git-like abstractions to avoid a high cognitive load.

Currently, those are:

- **Change:** a change of value on a column of a table. It can be a deletion or addition.

- **Commit:** a group of changes signed (ensuring persistence).

- **Branch:** a parallel context instantiation for communicating via multiple actors. It contains an index and multiple commits. For practical purposes, it stores all the credentials for the services it'll communicate.

- **Index:** it contains the group of changes that weren't committed and were done over the branch it belongs to.

- **Collaborator:** it's responsible for communicating with a specific service using its own interface. It can push, pull, delete and init.

- **Member:** a collaborator assigned to a table.

- **Team:** a group of members assigned to a schema.

- **Community:** a group of teams that are available to collaborate.

- **Project:** a group of schemas.

- **Owner:** it's responsible for orchestrating its own project given a community. It's a collaborator too.

- **Strategy:** strategy used to merge commits.

- **Pull request:** a group of commits performed by a team.
