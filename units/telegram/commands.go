package telegram

import (
	"context"
	"log"

	"bonefabric/adviser/store"
)

type command string

const commandStart command = "start"
const commandHelp command = "help"
const commandAddBookmark command = "addbookmark"

func (p *processor) cmdHelp() string {
	return messageHelp
}

func (p *processor) cmdAddBookmark(from int) string {
	p.setUserState(from, waitBookmarkName)
	return messageAddBookmark
}

func (p *processor) bookmarkNameReceived(name string, from int) string {
	p.setUserState(from, waitBookmarkText)
	u := p.state[from]
	u.meta.bookmark.name = name
	p.state[from] = u
	return messageAddBookmarkText
}

func (p *processor) bookmarkTextReceived(ctx context.Context, text string, from int) string {
	s := p.state[from]
	s.meta.bookmark.text = text

	p.resetUserState(from)

	err := p.store.Save(ctx, store.Bookmark{
		Name: s.meta.bookmark.name,
		Text: s.meta.bookmark.text,
		User: from,
	})

	if err != nil {
		log.Printf("failed to save bookmark: %s", err)
	}

	return messageBookmarkSaved
}
