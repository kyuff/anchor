# anchor

[![Build Status](https://github.com/kyuff/anchor/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kyuff/anchor/actions/workflows/go.yml)
[![Report Card](https://goreportcard.com/badge/github.com/kyuff/anchor)](https://goreportcard.com/report/github.com/kyuff/anchor/)
[![Go Reference](https://pkg.go.dev/badge/github.com/kyuff/anchor.svg)](https://pkg.go.dev/github.com/kyuff/anchor)
[![codecov](https://codecov.io/gh/kyuff/anchor/graph/badge.svg?token=GA4GSQCLZE)](https://codecov.io/gh/kyuff/anchor)

Library to manage application lifetime in a Go microservice architecture.

# Features

* Simple API to manage application lifetime
* Graceful shutdown of application components
* Freedom of choise for dependency injection
* Convenience methods to wire external APIs into Anchor

# Quickstart

cmd/main.go
```golang
package main

import (
    "github.com/kyuff/anchor"
    "example.com/myapp/internal/app""
)

func main() {
  os.Exit(app.Run(anchor.DefaultSignalWire()))
}
```

internal/app/anchor.go
```golang
package app

import (
    "github.com/kyuff/anchor"
)

func Run(wire anchor.Wire) int {
  return anchor.New(wire, anchor.WithDefaultSlog().
       Add(
         anchor.Close("database.Connection", func() io.Closer {
            return database()
         }),
         NewHttpServer()
       ).
       Run()
}
```

internal/app/database.go
```golang
package app

import (
    "database/sql"
    "os"
)

var database = anchor.Singleton(func() (*sql.DB, error) {
  return sql.Open("postgres", os.Getenv("DATABASE_URL")
})
```

internal/app/http.go
```golang
package app

import (
    "context"
    "net/http"
    "os"

    "example.com/myapp/internal/api"
)

type HttpServer struct {
    server *http.Server
}

func NewHttpServer() *HttpServer {
    return &HttpServer{
        server: &http.Server{
            Addr: os.Getenv("HTTP_ADDR"),
        },
    }
}

func (h *HttpServer) Setup(ctx context.Context) error {
    return api.Register(http.DefaultServeMux, database())
}

func (h *HttpServer) Start(ctx context.Context) error {
    return h.server.ListenAndServe()
}

func (h *HttpServer) Close(ctx context.Context) error {
    return h.server.Shutdown(ctx)
}
```
