package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte("{\"name\": \"golang-walletconnect-bridge\",\"description\": \"WalletConnect Bridge Server - GoLang\",\"version\": \"0.0.1-beta\"}"))
		if err != nil {
			return
		}
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
	fmt.Printf("Websocket Server listeing on addr: %s", *addr)
}
