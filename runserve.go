package main

import (
	"log"
	"net/http"
	//"fmt"
	//"os"
	//"html/template"
	"./shellscan"
	"./strscan"
)

func main() {
	http.HandleFunc("/strscan", strscan.Strscan)
	http.HandleFunc("/shellscan", shellscan.ShellScan)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}
