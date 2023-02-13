package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"
	"time"

	tgClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config"
	"bonefabric/adviser/config/yaml"
	"bonefabric/adviser/pool"
	"bonefabric/adviser/store"
	"bonefabric/adviser/store/mysql"
	"bonefabric/adviser/store/sqlite"
	tgUnit "bonefabric/adviser/units/telegram"
)

func main() {
	log.Println("application started")
	defer log.Println("application stopped")

	cnf, err := loadConfig()
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

	st, err := initStore(cnf, ctx)
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

func initStore(cnf config.Config, ctx context.Context) (store.Store, error) {
	ctx, canc := context.WithTimeout(ctx, time.Second*5)
	defer canc()

	switch cnf.StoreDriver() {
	case string(store.DriverSqlite3):
		return sqlite.New(cnf.StoreName())
	case string(store.DriverMysql):
		return mysql.New(ctx, mysql.DSN{
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

func loadConfig() (config.Config, error) {
	if flag.Parsed() {
		return nil, errors.New("failed to load config: flag already parsed")
	}
	confFile := flag.String("config", "config.yaml", "config file name")
	flag.Parse()

	if confFile == nil || *confFile == "" {
		return nil, errors.New("failed to load config: invalid file name")
	}

	switch filepath.Ext(*confFile) {
	case ".yaml":
		return yaml.Load(*confFile)
	default:
		return nil, errors.New("failed to load config: invalid file")
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
