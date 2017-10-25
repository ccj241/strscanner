package main

import (
  "net/http"
  "log"
  //"fmt"
  //"os"
  "./scan"
  //"html/template"
)

func main(){
    http.HandleFunc("/strscan",scan.Strscan)
    log.Fatal(http.ListenAndServe("localhost:8888",nil))
}
