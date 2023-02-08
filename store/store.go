package store

import "context"

type StoreDriver string

const StoreSqlite3 StoreDriver = "sqlite3"

type Store interface {
	Save(ctx context.Context, b Bookmark) error
	Delete(ctx context.Context, b Bookmark) error
	PickRandom(ctx context.Context, user int) (Bookmark, error)
	Exists(ctx context.Context, b Bookmark) (bool, error)
	Close() error
}

type Bookmark struct {
	Text string
	User int
}
