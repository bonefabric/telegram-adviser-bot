package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config"
	"bonefabric/adviser/config/yaml"
	"bonefabric/adviser/pool"
	"bonefabric/adviser/store"
	"bonefabric/adviser/store/mysql"
	"bonefabric/adviser/store/sqlite"
	tgUnit "bonefabric/adviser/units/telegram"
)

const configPath = "./config.yaml"

func main() {
	log.Println("application started")
	defer log.Println("application stopped")

	cnf, err := yaml.Load(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if cnf.Profiling() {
		defer profiler(cnf.ProfileFile())()
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleSysSignals(cancel)

	tgc := tgClient.New(cnf.TelegramToken())

	st, err := initStore(cnf)
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
}

func initStore(cnf config.Config) (store.Store, error) {
	switch cnf.StoreDriver() {
	case string(store.DriverSqlite3):
		return sqlite.New(cnf.StoreName())
	case string(store.DriverMysql):
		return mysql.New(mysql.DSN{
			UserName: cnf.StoreUser(),
			Password: cnf.StorePassword(),
			Host:     cnf.StoreHost(),
			Port:     cnf.StorePort(),
			DBName:   cnf.StoreName(),
		})
	default:
		return nil, errors.New("invalid driver")
	}
}

func handleSysSignals(call context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("os signal handled: %s\n", <-sig)
	call()
}

func profiler(fname string) func() {
	log.Println("profiling started")
	f, err := os.Create(fname)
	if err != nil {
		log.Printf("failed to open profiling file %s: %s\n", fname, err)
		return func() {}
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Printf("failed to start profiling: %s\n", err)
	}

	return func() {
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			log.Printf("failed to close profiling file: %s\n", err)
		}
		log.Println("profiler stopped")
	}
}
