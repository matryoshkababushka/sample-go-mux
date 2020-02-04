package main

import (
  "errors"
)

var (
  ErrNotInitialized = errors.New("Variable was not initialized")
)
