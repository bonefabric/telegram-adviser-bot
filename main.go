package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bonefabric/adviser/pool"
)

func main() {
	log.Println("application started")

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	p := pool.Pool{}
	p.Start(ctx)

	log.Println("application stopped")
}

func handleSysSignals(call context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("os signal handled: %s\n", <-sig)
	call()
}
