package telegram

import (
	"context"
	"regexp"
	"strings"

	"bonefabric/adviser/store"
)

type processor struct {
	state map[int]userState
	store store.Store
}

func (p *processor) process(ctx context.Context, msg string, from int) (string, error) {
	if strings.HasPrefix(msg, "/") {
		return p.processCmd(ctx, msg, from)
	}
	return p.processArg(ctx, msg, from)
}

// working with command
func (p *processor) processCmd(ctx context.Context, msg string, from int) (string, error) {
	p.resetUserState(from)

	switch p.extractCommand(msg) {
	case commandStart, commandHelp:
		return p.cmdHelp(), nil
	case commandAddBookmark:
		return p.cmdAddBookmark(from), nil
	case commandPickBookmark:
		return p.cmdPickBookmark(ctx, from), nil
	case commandRemoveBookmark:
		return p.cmdRemoveBookmark(from), nil
	default:
		return p.cmdHelp(), nil
	}
}

// working with argument
func (p *processor) processArg(ctx context.Context, arg string, from int) (string, error) {
	switch p.state[from].state {
	case defaultState:
		return p.cmdHelp(), nil
	case waitNewBookmarkName:
		return p.bookmarkNameReceived(arg, from), nil
	case waitNewBookmarkText:
		return p.bookmarkTextReceived(ctx, arg, from), nil
	case waitDeleteBookmarkName:
		return p.removingBookmarkNameReceived(ctx, from, arg), nil
	default:
		return p.cmdHelp(), nil
	}
}

func (p *processor) extractCommand(input string) command {
	re := regexp.MustCompile(`^/(\w+)`)
	match := re.FindStringSubmatch(input)
	if len(match) < 2 {
		return ""
	}
	return command(match[1])
}

func (p *processor) resetUserState(user int) {
	s := p.state[user]
	s.reset()
	p.state[user] = s
}

func (p *processor) setUserState(user int, st state) {
	s := p.state[user]
	s.state = st
	p.state[user] = s
}
