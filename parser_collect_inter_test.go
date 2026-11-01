package ansiterm

import "testing"

// collectInter used to append the intermediate byte to interBuffer but assign
// the result to paramBuffer, so interBuffer stayed permanently empty — no
// escape intermediate ever reached escDispatch — and the intermediate byte
// was deposited into the parameter buffer instead.

func TestEscapeIntermediateIsCollected(t *testing.T) {
	parser, _ := createTestParser("Escape")
	if _, err := parser.Parse([]byte("(")); err != nil {
		t.Fatal(err)
	}
	validateState(t, parser.currState, "EscapeIntermediate")
	if got := string(parser.context.interBuffer); got != "(" {
		t.Errorf("interBuffer = %q, want %q", got, "(")
	}
	if len(parser.context.paramBuffer) != 0 {
		t.Errorf("paramBuffer polluted with intermediate: %q", parser.context.paramBuffer)
	}
}

func TestIntermediatesAccumulate(t *testing.T) {
	parser, _ := createTestParser("Escape")
	if _, err := parser.Parse([]byte("( ")); err != nil {
		t.Fatal(err)
	}
	validateState(t, parser.currState, "EscapeIntermediate")
	if got := string(parser.context.interBuffer); got != "( " {
		t.Errorf("interBuffer = %q, want %q", got, "( ")
	}
	if len(parser.context.paramBuffer) != 0 {
		t.Errorf("paramBuffer polluted with intermediates: %q", parser.context.paramBuffer)
	}
}
