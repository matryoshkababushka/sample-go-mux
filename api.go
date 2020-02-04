package main

import (
  "github.com/crackhd/env"
  "log"
  "net/http"
	"github.com/gorilla/mux"
  "strings"
)

type API struct {
  env       * env.Env
  backend   BackendInterface
  imageURL  string
}

func NewAPI(env *env.Env, backend BackendInterface) *API {
  if backend == nil {
    panic(ErrNotInitialized)
  }

  imageURL, ok := env.Get("IMAGE_URL", "http://127.0.0.1:8080/uploads/")
  if !ok {
    log.Println("WARNING: IMAGE_URL variable is not set, image upload won't work properly.")
  }
  imageURL = strings.TrimRight(imageURL, "/")
  return &API{env, backend, imageURL}
}

func (api * API) Engine() *Engine {
  return api.backend.Engine()
}

func (api * API) Run() error {
  r := mux.NewRouter()

  b := api.backend
  e := b.Engine()

  r.HandleFunc("/", api.IndexController).Methods("GET")

	r.HandleFunc("/auth", e.AuthController)

  r.HandleFunc("/example", api.ExampleController).Methods("GET")

  // Static file serving. In production, should be overriden by nginx.
  NeuteredFileServer(NeuteredConfig{
    DirectoryListings:  false,        // Switches implementation
    Path:               "/uploads",   // HTTP/URL path prefix
    Directory:          "./public/",  // Filesystem path prefix
    Methods:            "GET",        // Allowed methods (??)
    Router:             r,            // Apply to router
  })

  h := e.CorsHandler()

  log.Println("Server is running on port 8080")
  return http.ListenAndServe(":8080", h(r))
}

func (api * API) IndexController(w http.ResponseWriter, r *http.Request) {
  b := api.backend
  e := b.Engine()
  s := e.Session(r, w)

  s = s // TODO: Magic?

  w.WriteHeader(b.ExampleGetStatus())
  w.Write([]byte("Hello world!"))
}

func (api * API) ExampleController(w http.ResponseWriter, r *http.Request) {
  b := api.backend
  e := b.Engine()
  s := e.Session(r, w)

  status := http.StatusInternalServerError

  s = s // TODO: magic

  w.WriteHeader(status)
}
