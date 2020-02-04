package main

import (
  "github.com/crackhd/env"
  "net/http"
)

type BackendInterface interface {
  Close()
  Engine() *Engine
  ExampleGetStatus() int
}

type Backend struct {
  engine    * Engine
}

func NewBackendBase(env *env.Env) BackendInterface {
  e := NewEngine(env)
  return &Backend{e}
}

func (b *Backend) Close() {

}

func (b *Backend) ExampleGetStatus() int {
  return http.StatusInternalServerError
}

func (b *Backend) Engine() *Engine {
  return b.engine
}
