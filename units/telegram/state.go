package telegram

type state int

type userState struct {
	state state
	meta  userMeta
}

type userMeta struct {
	bookmark struct {
		name string
		text string
	}
}

const (
	defaultState state = iota
	waitNewBookmarkName
	waitNewBookmarkText
	waitDeleteBookmarkName
)

func (s *userState) reset() {
	s.state = defaultState
	s.meta = userMeta{}
}
