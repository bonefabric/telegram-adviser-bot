package telegram

import (
	"context"
	"log"
	"sync"
	"time"

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

	done := make(chan struct{})
	go t.start(ctx, done)

	<-done
	log.Println("telegram unit stopped")
	wg.Done()
}

func (t *Telegram) start(ctx context.Context, done chan<- struct{}) {
	defer close(done)

cycle:
	for {
		select {
		case <-ctx.Done():
			break cycle
		default:
			time.Sleep(time.Second)
			c, canc := context.WithTimeout(ctx, time.Second*5)
			updates := t.fetch(c)
			canc()

			for _, u := range updates {
				c, canc = context.WithTimeout(ctx, time.Second*3)
				t.process(c, u)
				canc()
			}
		}
	}
	done <- struct{}{}
}

func (t *Telegram) fetch(ctx context.Context) []telegram.Update {
	updates, err := t.client.Updates(ctx)
	if err != nil {
		log.Printf("failed to fetch updates: %s", err)
		return nil
	}
	return updates
}

func (t *Telegram) process(ctx context.Context, upd telegram.Update) {
	//todo
	log.Println("upd")
}
