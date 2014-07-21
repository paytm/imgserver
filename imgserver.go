package main

import (
  "net/http"
  "github.com/codegangsta/negroni"
  "code.google.com/p/gcfg"
  "catalog"
  "flag"
  "log"
)

type Config struct {
  DB struct {
    DSN string
  }
  Server struct {
    Port string
  }
}

func readConfig(cfg *Config,path string) bool {
  err := gcfg.ReadFileInto(cfg,path + "/imgserver.ini")
  if err == nil {
    log.Println("read config from ",path)
    return true
  }
  return false
}

func main() {

  var cfg Config

  cfg.Server.Port = "9999" // default port

  ok := readConfig(&cfg, ".") || readConfig(&cfg,"/etc")
  if !ok {
     log.Fatal("failed to read config")
  }

  flag.Parse()

  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    http.Error(w, "File not found", http.StatusNotFound)
  })

  n := negroni.Classic()
  n.Use(negroni.HandlerFunc(catalog.ImageRedir(cfg.DB.DSN)));
  n.UseHandler(mux)
  n.Run(":" + cfg.Server.Port)
}
