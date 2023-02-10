package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config/yaml"
	"bonefabric/adviser/pool"
	"bonefabric/adviser/store"
	"bonefabric/adviser/store/sqlite"
	tgUnit "bonefabric/adviser/units/telegram"

	_ "github.com/mattn/go-sqlite3"
)

const configPath = "./config.yaml"

func main() {
	log.Println("application started")

	cnf, err := yaml.Load(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	tgc := tgClient.New(cnf.TelegramToken())

	st, err := initStore(cnf.StoreDriver())
	if err != nil {
		log.Fatalf("failed to init store: %s", err)
	}

	defer func(st store.Store) {
		if err := st.Close(); err != nil {
			log.Printf("failed to close store: %s", err)
		}
	}(st)

	tgu := tgUnit.New(tgc, st)

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

func initStore(driver string) (store.Store, error) {
	switch driver {
	case string(store.DriverSqlite3):
		return sqlite.New("data")
	default:
		return nil, errors.New("invalid driver")
	}
}
