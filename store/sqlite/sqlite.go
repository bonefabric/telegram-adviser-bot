package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"bonefabric/adviser/store"
)

const initial = `CREATE TABLE IF NOT EXISTS 
    bookmarks (
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(255) DEFAULT NULL,
    text TEXT DEFAULT NULL,
    user INT
    )`

type Sqlite struct {
	db *sql.DB
}

func New(fname string) (*Sqlite, error) {
	var s Sqlite

	db, err := sql.Open("sqlite3", fname)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initial)
	if err != nil {
		return nil, fmt.Errorf("failed to init sqlite3 store: %s", err)
	}

	s.db = db
	return &s, nil
}

func (s Sqlite) Save(ctx context.Context, b store.Bookmark) error {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO bookmarks (name, text, user) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to store bookmark: %s", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.Exec(b.Name, b.Text, b.User)
	return err
}

func (s Sqlite) Delete(ctx context.Context, b store.Bookmark) error {
	stmt, err := s.db.PrepareContext(ctx, `DELETE FROM bookmarks WHERE name = ? AND user = ?`)
	if err != nil {
		return fmt.Errorf("failed to remove bookmark: %s", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.Exec(b.Name, b.User)
	return err
}

func (s Sqlite) PickRandom(ctx context.Context, user int) (store.Bookmark, error) {
	var b store.Bookmark
	b.User = user

	stmt, err := s.db.PrepareContext(ctx, `SELECT name, text FROM bookmarks WHERE user = ? ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		return b, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if err = stmt.QueryRow(user).Scan(&b.Name, &b.Text); err != nil {
		return b, err
	}

	return b, nil
}

func (s Sqlite) Exists(ctx context.Context, b store.Bookmark) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s Sqlite) Close() error {
	return s.db.Close()
}
