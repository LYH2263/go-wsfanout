package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/wsfanout"
	"example.com/wsfanout/internal/api"
)

func main() {
	addr := flag.String("addr", ":8102", "listen address")
	web := flag.String("web", "web", "static web directory")
	flag.Parse()

	hub := wsfanout.New()
	defer hub.Close()

	srv := api.New(hub, *web)
	go func() {
		log.Printf("wsd listening on %s", *addr)
		if err := srv.ListenAndServe(*addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("shutting down")
}
