package store

import "context"

type Store interface {
	Save(ctx context.Context, b Bookmark) error
	Delete(ctx context.Context, b Bookmark) error
	PickRandom(ctx context.Context) ([]Bookmark, error)
	Exists(ctx context.Context) (bool, error)
}

type Bookmark struct {
	Text string
	User int
}
