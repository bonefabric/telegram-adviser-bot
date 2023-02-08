package pool

import (
	"context"
	"sync"
)

type Pool struct {
	units []Unit
}

func (p *Pool) AddUnits(units ...Unit) {
	p.units = append(p.units, units...)
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(p.units))
	for _, u := range p.units {
		u.Start(ctx, &wg)
	}
	wg.Wait()
}
