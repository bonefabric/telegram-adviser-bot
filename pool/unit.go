package pool

import (
	"context"
	"sync"
)

type Unit interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}
