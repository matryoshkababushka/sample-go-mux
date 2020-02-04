package main

import (
    "net/http"
    "strings"
    "github.com/gorilla/mux"
)

// Fix HTTP static server lists directory contents.
// Works both with gorilla/mux and net/http
// This implementation only serves files for security.
// Kudos Alex Edwards
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings

type NeuteredConfig struct {
  DirectoryListings         bool
  Path, Directory, Methods  string
  Router                    *mux.Router     // For gorilla/mux
  Server                    *http.ServeMux  // OR net/http
}

func NeuteredFileServer(config NeuteredConfig) {

  if config.Path != "/" {
    config.Path = strings.TrimRight(config.Path, "/")
  }

  var fileServer http.Handler
  if config.DirectoryListings {
    fileServer = http.FileServer(http.Dir(config.Directory))
  } else {
    fileServer = http.FileServer(neuteredFileSystem{http.Dir(config.Directory)})
  }

  if config.Router != nil {
    config.Router.Handle(config.Path, http.NotFoundHandler())
    config.Router.PathPrefix(config.Path + "/").Handler(http.StripPrefix(config.Path, fileServer))
  } else if config.Server != nil {
    config.Server.Handle(config.Path, http.NotFoundHandler())
    config.Server.Handle(config.Path + "/", http.StripPrefix("/static", fileServer))
  } else {
    panic("Use Router or Server")
  }
}

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := strings.TrimSuffix(path, "/") + "/index.html"
        if _, err := nfs.fs.Open(index); err != nil {
            return nil, err
        }
    }

    return f, nil
}
