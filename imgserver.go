package main

import (
  "net/http"
  "github.com/codegangsta/negroni"
  "catalog"
  "flag"
)

func main() {

  flag.Parse()

  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    http.Error(w, "File not found", http.StatusNotFound)
  })

  n := negroni.Classic()
  n.Use(negroni.HandlerFunc(catalog.ImageRedir()));
  n.UseHandler(mux)
  n.Run(":9999")
}
