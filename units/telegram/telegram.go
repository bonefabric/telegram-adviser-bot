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
	defer func(d chan<- struct{}) {
		d <- struct{}{}
		close(d)
	}(done)

	us := make(chan telegram.Update, 100)
	usdone := make(chan struct{})

	go t.process(us, usdone)

out:
	for {
		select {
		case <-ctx.Done():
			close(us)
			<-usdone
			break out
		default:
			time.Sleep(time.Second * 3)
			c, canc := context.WithTimeout(ctx, time.Second*5)
			updates := t.fetch(c)
			canc()
			for _, u := range updates {
				us <- u
			}
		}
	}
}

func (t *Telegram) fetch(ctx context.Context) []telegram.Update {
	updates, err := t.client.Updates(ctx)
	if err != nil {
		log.Printf("failed to fetch updates: %s", err)
		return nil
	}
	return updates
}

func (t *Telegram) process(upds <-chan telegram.Update, done chan<- struct{}) {
	defer func(d chan<- struct{}) {
		d <- struct{}{}
		close(d)
	}(done)

	for {
		upd, ok := <-upds
		if !ok {
			break
		}

		msg := upd.Message
		if msg == nil || msg.Text == nil || *msg.Text == "" || msg.From == nil {
			log.Printf("failed to process update %d: empty message, text or user\n", upd.ID)
			continue
		}

		b := store.Bookmark{
			Text: *msg.Text,
			User: msg.From.ID,
		}

		ctx, canc := context.WithTimeout(context.Background(), time.Second)
		saveErr := t.store.Save(ctx, b)
		canc()

		if saveErr != nil {
			log.Printf("failed to store bookmark: %s\n", saveErr)
		} else {
			log.Printf("bookmark stored from user id %d\n", b.User)
		}
	}
}
