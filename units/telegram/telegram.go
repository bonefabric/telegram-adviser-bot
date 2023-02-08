package telegram

import (
	"context"
	"log"
	"sync"
	"time"

	"bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/store"
)

type Telegram struct {
	client telegram.Telegram
	store  store.Store
}

func New(client telegram.Telegram, store store.Store) Telegram {
	return Telegram{
		client: client,
		store:  store,
	}
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
	msg := upd.Message
	if msg == nil || msg.Text == nil || *msg.Text == "" || msg.From == nil {
		log.Printf("failed to process update %d: empty message, text or user\n", upd.ID)
		return
	}

	b := store.Bookmark{
		Text: *msg.Text,
		User: msg.From.ID,
	}
	if err := t.store.Save(ctx, b); err != nil {
		log.Printf("failed to store bookmark: %s\n", err)
	}
}
