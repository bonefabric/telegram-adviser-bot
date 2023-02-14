package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	telegramClient "bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/config"
	"bonefabric/adviser/config/yaml"
	"bonefabric/adviser/pool"
	"bonefabric/adviser/store"
	"bonefabric/adviser/store/mysql"
	"bonefabric/adviser/store/sqlite"
	"bonefabric/adviser/units/telegram"
)

const (
	defaultConfigFile = "config.yaml"
	maxStoreInitTime  = time.Second * 3
)

type App struct {
	initialized bool
	config      config.Config
	context     struct {
		ctx    context.Context
		cancel context.CancelFunc
	}
	profiling struct {
		enabled      bool
		profilerFunc func()
	}
	pool  pool.Pool
	store store.Store
}

func New() App {
	return App{}
}

func (a *App) Init() error {
	log.Println("app initialization started")

	confFile, err := a.parseConfigFlag()
	if err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	a.config, err = a.loadConfig(confFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err)
	}

	a.context.ctx, a.context.cancel = context.WithCancel(context.Background())
	go a.handleSysSignals()

	if a.store, err = a.initStore(); err != nil {
		return fmt.Errorf("failed to init store: %s", err)
	}
	a.fillPool()

	a.initialized = true
	log.Println("app initialization finished successful")
	return nil
}

func (a *App) Run() {
	log.Println("app started")
	defer a.context.cancel()

	defer func(store store.Store) {
		if err := store.Close(); err != nil {
			log.Printf("failed to close store: %s\n", err)
		}
	}(a.store)

	a.pool.Start(a.context.ctx)
	log.Println("app finished")
}

func (a *App) parseConfigFlag() (cnf string, err error) {
	if flag.Parsed() {
		err = errors.New("flag already parsed")
		return
	}
	flag.StringVar(&cnf, "config", defaultConfigFile, "config file")
	flag.Parse()
	return
}

func (a *App) loadConfig(name string) (config.Config, error) {
	switch filepath.Ext(name) {
	case ".yaml":
		return yaml.Load(name)
	default:
		return nil, errors.New("invalid config file extension")
	}
}

func (a *App) handleSysSignals() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("os signal handled: %s\n", <-sig)
	a.context.cancel()
}

func (a *App) initStore() (store.Store, error) {
	ctx, c := context.WithTimeout(a.context.ctx, maxStoreInitTime)
	defer c()

	switch a.config.StoreDriver() {
	case string(store.DriverMysql):
		s, err := mysql.New(ctx, mysql.DSN{
			UserName: a.config.StoreUser(),
			Password: a.config.StorePassword(),
			Host:     a.config.StoreHost(),
			Port:     a.config.StorePort(),
			DBName:   a.config.StoreName(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init mysql store: %s", err)
		}
		return s, nil
	case string(store.DriverSqlite3):
		s, err := sqlite.New(a.config.StoreName())
		if err != nil {
			return nil, fmt.Errorf("failed to init sqlite3 store: %s", err)
		}
		return s, nil
	default:
		return nil, errors.New("invalid driver")
	}
}

func (a *App) fillPool() {
	telegramUnit := telegram.New(telegramClient.New(a.config.TelegramToken()), a.store)
	a.pool.AddUnits(&telegramUnit)
}
