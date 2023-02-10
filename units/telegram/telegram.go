package telegram

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"bonefabric/adviser/clients/telegram"
	"bonefabric/adviser/store"
)

const onErrAnswer = "Sorry, something went wrong. Try again later"

// Telegram unit
type Telegram struct {
	client    telegram.Telegram
	store     store.Store
	processor processor
}

// New Telegram unit constructor
func New(client telegram.Telegram, store store.Store) Telegram {
	return Telegram{
		client:    client,
		store:     store,
		processor: processor{},
	}
}

// Start main Telegram working process
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

	go t.handleUpdates(us, usdone)

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
			updates := t.fetchUpdates(c)
			canc()
			for _, u := range updates {
				us <- u
			}
		}
	}
}

func (t *Telegram) fetchUpdates(ctx context.Context) []telegram.Update {
	updates, err := t.client.Updates(ctx)
	if err != nil {
		log.Printf("failed to fetchUpdates updates: %s", err)
		return nil
	}
	return updates
}

func (t *Telegram) handleUpdates(upds <-chan telegram.Update, done chan<- struct{}) {
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
			log.Printf("failed to handleUpdates update %d: empty message, text or user\n", upd.ID)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		answer, err := t.processor.process(ctx, *msg.Text)
		if err != nil {
			log.Printf("failed to process command: %s\n", err)
			answer = onErrAnswer
		}
		if err = t.sendAnswer(ctx, answer, msg.From.ID); err != nil {
			log.Printf("failed to process command: %s\n", err)
		}
		cancel()
	}
}

func (t *Telegram) sendAnswer(ctx context.Context, answer string, user int) error {
	err := t.client.SendMessage(ctx, telegram.SendMessageOptions{
		ChatID: user,
		Text:   answer,
	})
	if err != nil {
		return fmt.Errorf("faield to send answer: %s", err)
	}
	return nil
}
