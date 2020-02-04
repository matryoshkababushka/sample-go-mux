package main

import (
  "net/http"
	"github.com/gorilla/sessions"
  "github.com/crackhd/env"
  "log"
  "strings"
  "time"
  "github.com/rs/cors"
)

// TODO: encapsulate within Engine{}
var (
  SITE_URL      string
  SITE_NAME     string
  COOKIE_DOMAIN string
)

const EMAIL_INTERVAL    int64 = 60
const SESSION_TIMEOUT   int64 = 60 * 60 * 2

type Engine struct {
  store       *sessions.CookieStore
  sessionKey  string
}

func NewEngine(env * env.Env) *Engine {
  COOKIE_DOMAIN, _ = env.Get("COOKIE_DOMAIN", "localhost")
  SITE_URL, _ = env.Get("SITE_URL", "http://localhost:3000/")
  SITE_NAME, _ = env.Get("SITE_NAME", "Changeme")

  sessionKey, ok := env.Get("SESSION_KEY", "changeme")
  if !ok {
    log.Println("WARNING! SESSION_KEY secret is not set!")
  }

  store := sessions.NewCookieStore([]byte(sessionKey))

  return &Engine{store, sessionKey}
}

func (engine Engine) Session(r *http.Request, w http.ResponseWriter) * sessions.Session {
  s, err := engine.store.Get(r, "auth")
  if err != nil {
    log.Println("Failed to get s for request, error:", err.Error())
    return nil
  }
  s.Options = &sessions.Options{
    Path:     "/",
    Domain:   COOKIE_DOMAIN,
    MaxAge:   int(SESSION_TIMEOUT),
    Secure:   !engine.IsLocalhostDevMode(),
    HttpOnly: true,
  }
  return s
}

func (engine Engine) Save(s *sessions.Session, r *http.Request, w http.ResponseWriter) error {
  if s == nil {
    panic(ErrNotInitialized)
  }

  s.Values["timestamp"] = time.Now().Unix()
  return s.Save(r, w)
}

func (engine Engine) IsLocalhostDevMode() bool {
  return strings.HasPrefix(COOKIE_DOMAIN, "localhost")
}

func (engine Engine) GetUserURL() string {
  return SITE_URL
}

func (engine Engine) CorsHandler() func(http.Handler) http.Handler {
  origins := []string{strings.TrimRight(engine.GetUserURL(), "/")}

  methods := []string{
    http.MethodOptions,
    http.MethodGet,
    http.MethodPost,
  }

  headers := []string{"Cookie","Content-Type", "X-Requested-With"}

  corsOpts := cors.New(cors.Options{
    AllowedOrigins:   origins,
    AllowedMethods:   methods,
    AllowedHeaders:   headers,
    AllowCredentials: true,
  })

  log.Println("[CORS] Allow origins:", origins)

  return corsOpts.Handler

  // This one just did not work ("github.com/gorilla/handlers") --
  //corsOrigins := handlers.AllowedOrigins(origins)
  //corsMethods := handlers.AllowedMethods(methods)
  //corsCredentials := handlers.AllowCredentials()
  //corsHeaders := handlers.AllowedHeaders(headers)
  //return handlers.CORS(corsOrigins, corsMethods, corsCredentials, corsHeaders)
}

func (engine Engine) AuthController(w http.ResponseWriter, r *http.Request) {

  s := engine.Session(r, w)

  status := http.StatusInternalServerError

  s = s // TODO: magic

  w.WriteHeader(status)
}
