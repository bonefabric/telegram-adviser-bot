package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config"
	"bonefabric/adviser/pool"
	tgUnit "bonefabric/adviser/units/telegram"
)

func main() {
	log.Println("application started")

	cnf := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	tgc := tgClient.New(cnf.TgToken())

	tgu := tgUnit.New(tgc)

	p := pool.Pool{}
	p.AddUnits(&tgu)
	p.Start(ctx)

	log.Println("application stopped")
}

func handleSysSignals(call context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("os signal handled: %s\n", <-sig)
	call()
}
