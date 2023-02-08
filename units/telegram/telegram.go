package telegram

import (
	"context"
	"log"
	"sync"
	"time"
)

type Telegram struct {
}

func (t Telegram) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("telegram unit started")
	<-ctx.Done()
	time.Sleep(time.Second)
	log.Println("telegram unit started")
	wg.Done()
}
