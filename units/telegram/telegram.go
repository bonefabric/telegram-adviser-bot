package telegram

import (
	"context"
	"log"
	"sync"

	"bonefabric/adviser/clients/telegram"
)

type Telegram struct {
	client telegram.Telegram
}

func New(client telegram.Telegram) Telegram {
	return Telegram{client: client}
}

func (t *Telegram) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("telegram unit started")

	errs := make(chan error)
	go t.start(ctx, errs)

	if err := <-errs; err != nil {
		log.Printf("telegram unit stopped with error: %s\n", err)
	} else {
		log.Println("telegram unit stopped")
	}
	wg.Done()
}

func (t *Telegram) start(ctx context.Context, errs chan<- error) {
	defer close(errs)

	//todo
	<-ctx.Done()
	errs <- nil
}
