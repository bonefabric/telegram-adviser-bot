package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"bonefabric/adviser/store"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db *sql.DB
}

type DSN struct {
	UserName string
	Password string
	Host     string
	Port     int
	DBName   string
}

const initial = `
CREATE TABLE IF NOT EXISTS bookmarks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  text TEXT NOT NULL,
  user INT NOT NULL
);
`

func New(ctx context.Context, dsn DSN) (*Mysql, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dsn.UserName, dsn.Password, dsn.Host, dsn.Port, dsn.DBName))
	if err != nil {
		return nil, err
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	if _, err = db.Exec(initial); err != nil {
		return nil, err
	}
	return &Mysql{db: db}, nil
}

func (s *Mysql) Save(ctx context.Context, b store.Bookmark) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO bookmarks (name, text, user) VALUES (?, ?, ?)",
		b.Name, b.Text, b.User)
	return err
}

func (s *Mysql) Delete(ctx context.Context, b store.Bookmark) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM bookmarks WHERE name = ? AND user = ?", b.Name, b.User)
	return err
}

func (s *Mysql) PickRandom(ctx context.Context, user int) (store.Bookmark, error) {
	var b store.Bookmark
	err := s.db.QueryRowContext(ctx,
		"SELECT name, text FROM bookmarks WHERE user = ? ORDER BY RANDOM() LIMIT 1", user).Scan(&b.Name, &b.Text)
	if err != nil {
		return b, err
	}
	b.User = user
	return b, nil
}

func (s *Mysql) Exists(ctx context.Context, b store.Bookmark) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bookmarks WHERE name = ? AND user = ?",
		b.Name, b.User).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Mysql) Close() error {
	return s.db.Close()
}
