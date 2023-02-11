package store

import "context"

type Driver string

const DriverSqlite3 Driver = "sqlite3"

type Store interface {
	Save(ctx context.Context, b Bookmark) error
	Delete(ctx context.Context, b Bookmark) error
	PickRandom(ctx context.Context, user int) (Bookmark, error)
	Exists(ctx context.Context, b Bookmark) (bool, error)
	Close() error
}

type Bookmark struct {
	Name string
	Text string
	User int
}
