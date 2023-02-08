package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config"
	"bonefabric/adviser/pool"
	"bonefabric/adviser/store"
	"bonefabric/adviser/store/sqlite"
	tgUnit "bonefabric/adviser/units/telegram"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("application started")

	cnf := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	tgc := tgClient.New(cnf.TgToken())

	st, err := initStore(cnf.StoreDriver())
	if err != nil {
		log.Fatalf("failed to init store: %s", err)
	}

	defer func(st store.Store) {
		if err := st.Close(); err != nil {
			log.Printf("failed to close store: %s", err)
		}
	}(st)

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

func initStore(driver store.StoreDriver) (store.Store, error) {
	switch driver {
	case store.StoreSqlite3:
		return sqlite.New("data")
	default:
		return nil, errors.New("invalid driver")
	}
}
