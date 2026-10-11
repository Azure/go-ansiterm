package ansiterm

type oscStringState struct {
	baseState
}

func (oscState oscStringState) Handle(b byte) (s state, e error) {
	oscState.parser.logf("OscString::Handle %#x", b)

	// C1 controls use bytes 0x80-0x9F. Those bytes are also UTF-8
	// continuations. U+0410 is D0 90, and 0x90 is DCS, so a title
	// containing that character left the string and the next byte
	// aborted Parse. A continuation that belongs to a sequence started
	// in this string stays here. A lone 0x90 still enters DCS, and a
	// lone 0x9C still ends the string.
	if oscState.parser.oscUtf8Remain > 0 && b >= 0x80 && b <= 0xBF {
		oscState.parser.oscUtf8Remain--
		return oscState, nil
	}
	oscState.parser.oscUtf8Remain = utf8Remain(b)

	nextState, err := oscState.baseState.Handle(b)
	if nextState != nil || err != nil {
		oscState.parser.oscUtf8Remain = 0
		return nextState, err
	}

	// There are several control characters and sequences which can
	// terminate an OSC string. Most of them are handled by the baseState
	// handler. The ANSI_BEL character is a special case which behaves as a
	// terminator only for an OSC string.
	if b == ANSI_BEL {
		oscState.parser.oscUtf8Remain = 0
		return oscState.parser.ground, nil
	}

	return oscState, nil
}

// utf8Remain reports how many continuation bytes follow b.
func utf8Remain(b byte) int {
	switch {
	case b >= 0xC2 && b <= 0xDF:
		return 1
	case b >= 0xE0 && b <= 0xEF:
		return 2
	case b >= 0xF0 && b <= 0xF4:
		return 3
	default:
		return 0
	}
}
