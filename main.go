package main

import (
  "log"
  "github.com/crackhd/env"
)

var (
  e env.Env
)

func init() {
  var err error

	if err = e.LoadFile(); err != nil {
		log.Println("WARN: Error reading .env file:", err.Error())
	}
}

func main() {
  log.Println("changeme works!")

  b := NewBackendBase(&e)
  defer b.Close()

  a := NewAPI(&e, b)
  a.Run()
}
