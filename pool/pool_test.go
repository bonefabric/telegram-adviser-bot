package pool

import (
	"context"
	"sync"
	"testing"
)

type mockUnit struct {
	startCalled bool
}

func (u *mockUnit) Start(_ context.Context, wg *sync.WaitGroup) {
	u.startCalled = true
	wg.Done()
}

func TestPool_Start(t *testing.T) {
	t.Run("starts all units", func(t *testing.T) {
		ctx := context.Background()
		p := &Pool{}
		u1 := &mockUnit{}
		u2 := &mockUnit{}
		p.AddUnits(u1, u2)

		p.Start(ctx)

		if !u1.startCalled {
			t.Error("unit 1 was not started")
		}
		if !u2.startCalled {
			t.Error("unit 2 was not started")
		}
	})

	t.Run("waits for all units to finish", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		p := &Pool{}
		u1 := &mockUnit{}
		p.AddUnits(u1)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			p.Start(ctx)
			wg.Done()
		}()
		cancel()
		wg.Wait()

		if !u1.startCalled {
			t.Error("unit 1 was not started")
		}
	})
}
