# TODOs: A todo-list written in Go, with Clean Architecture.


[![Build Status](https://drone.unknowntpo.net/api/badges/unknowntpo/todos/status.svg?ref=refs/heads/master)](https://drone.unknowntpo.net/unknowntpo/todos)
[![Go.Dev](https://godoc.org/github.com/unknowntpo/todos?status.svg=)](https://pkg.go.dev/github.com/unknowntpo/todos?utm_source=godoc) [![Go Report Card](https://goreportcard.com/badge/github.com/unknowntpo/todos)](https://goreportcard.com/report/github.com/unknowntpo/todos) [![codecov](https://codecov.io/gh/unknowntpo/todos/branch/master/graph/badge.svg?token=UV6IIUUCW2)](https://codecov.io/gh/unknowntpo/todos)

> If github failed to render mermaid graphs, click the button below to see the Documentation at HackMD
> 
[![hackmd-github-sync-badge](https://hackmd.io/niPMhxhbSg-rNNzsj0Hrsw/badge)](https://hackmd.io/niPMhxhbSg-rNNzsj0Hrsw)

> It's a project based on [Advanced patterns for building APIs and web applications in Go](https://lets-go-further.alexedwards.net/),

> I modify the architecture of code to clean architecture, change the purpose of the code (from movie information service to todo list),  add some tests, deploy it to my own server in digitalocean.
> You can visit the api at [API Endpoint](https://todos.unknowntpo.net/v1/healthcheck)



[TOC]
## Quick Start
Clone the repository
```
$  git clone https://github.com/unknowntpo/todos.git
```

Use the default `.envrc` file
```
cp .envrc.example .envrc
```

Set up [mailtrap](https://mailtrap.io) for receiving email (for testing purpose)

Change username and password in `app_config-dev.yml` to your username and password.

```
smtp:
    host: "smtp.mailtrap.io"
    port: 25
    username: "<your-username>"
    password: "<your-password>"
    sender: "TODOs <no-reply@todos.unknowntpo.net>"
```

Run the project
```
$ make run/compose/up
```

Stop the project
```
$ make run/compose/down
```
## Project Walkthrough
[TODOs - Project Walkthrough](https://hackmd.io/@unknowntpo/todos-project-walkthrough)

:construction: Not Finished!
