package telegram

type command string

const commandHelp command = "help"

func (p *processor) cmdHelp() string {
	return helpMessage
}
