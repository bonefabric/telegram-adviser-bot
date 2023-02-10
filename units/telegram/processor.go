package telegram

import (
	"context"
	"regexp"
)

type processor struct {
}

func (p *processor) extractCommand(input string) command {
	re := regexp.MustCompile(`^/(\w+)`)
	match := re.FindStringSubmatch(input)
	if len(match) < 2 {
		return ""
	}
	return command(match[1])
}

func (p *processor) process(ctx context.Context, msg string) (string, error) {
	switch p.extractCommand(msg) {
	case commandHelp:
		return p.cmdHelp(), nil
	default:
		return p.cmdHelp(), nil
	}
}
