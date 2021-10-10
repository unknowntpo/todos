# TODOs: A todo-list written in Go, with Clean Architecture.

[![hackmd-github-sync-badge](https://hackmd.io/niPMhxhbSg-rNNzsj0Hrsw/badge)](https://hackmd.io/niPMhxhbSg-rNNzsj0Hrsw)

[![Build Status](https://cloud.drone.io/api/badges/unknowntpo/todos/status.svg)](https://cloud.drone.io/unknowntpo/todos) [![Go.Dev](https://godoc.org/github.com/unknowntpo/todos?status.svg=)](https://pkg.go.dev/github.com/unknowntpo/todos?utm_source=godoc) [![Go Report Card](https://goreportcard.com/badge/github.com/unknowntpo/todos)](https://goreportcard.com/report/github.com/unknowntpo/todos) [![codecov](https://codecov.io/gh/unknowntpo/todos/branch/master/graph/badge.svg?token=UV6IIUUCW2)](https://codecov.io/gh/unknowntpo/todos)

[TOC]
## Quick Start

## Project Walkthrough
### Database Schema
```mermaid
erDiagram
          User ||..o{ Token: has
          User ||..o{ Task: has
```
### Clean Architecture
### Makefile
### Testing
* Testcontainers v.s. sqlmock
    * Test on real database
* use build tag to seperate integration test and unit test

### Dockerfile
#### Multi-stage build

```mermaid
graph TD
    A(config-base) -->            
    | copy Makefile<br/>.envrc<br/>golang-migrate binary file<br/>migration files<br/>testdata<br/>config.sh|B(config)
    C(build-base<br/><build the binary file>) -->
    |copy binary file<br/>app_config-prod.yml| F(production)
    E(scratch) -->|as base image| F(production)
```
* Parallel build
    * When we change content in build base, config won't be changed.
* `.envrc`
### Configuration management
### Graceful shutdown