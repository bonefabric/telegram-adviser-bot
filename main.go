package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/pool"
	tgUnit "bonefabric/adviser/units/telegram"
)

func main() {
	log.Println("application started")

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	tg := tgUnit.New(tgClient.Telegram{})

	p := pool.Pool{}
	p.AddUnits(&tg)
	p.Start(ctx)

	log.Println("application stopped")
}

func handleSysSignals(call context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("os signal handled: %s\n", <-sig)
	call()
}
