package ansiterm

import (
	"fmt"
	"testing"
)

func getStateNames() []string {
	parser, _ := createTestParser("Ground")

	var stateNames []string
	for _, s := range parser.stateMap {
		stateNames = append(stateNames, s.Name())
	}

	return stateNames
}

func stateTransitionHelper(t *testing.T, start string, end string, bytes []byte) {
	t.Helper()
	for _, b := range bytes {
		t.Run(fmt.Sprintf("Start=%s/%q", start, b), func(t *testing.T) {
			parser, _ := createTestParser(start)
			if _, err := parser.Parse([]byte{b}); err != nil {
				t.Fatal(err)
			}
			validateState(t, parser.currState, end)
		})
	}
}

func anyToXHelper(t *testing.T, bytes []byte, expectedState string) {
	t.Helper()
	for _, s := range getStateNames() {
		stateTransitionHelper(t, s, expectedState, bytes)
	}
}

func funcCallParamHelper(t *testing.T, bytes []byte, start string, expected string, expectedCalls []string) {
	t.Helper()
	t.Run(fmt.Sprintf("Start=%s/%q", start, bytes), func(t *testing.T) {
		t.Helper()
		parser, evtHandler := createTestParser(start)
		if _, err := parser.Parse(bytes); err != nil {
			t.Errorf("Error parsing %q: %v", string(bytes), err)
		}
		validateState(t, parser.currState, expected)
		validateFuncCalls(t, evtHandler.FunctionCalls, expectedCalls)
	})
}

func parseParamsHelper(t *testing.T, bytes []byte, expectedParams []string) {
	t.Helper()
	params := parseParams(bytes)

	if len(params) != len(expectedParams) {
		t.Errorf("Parsed   parameters: %v", params)
		t.Errorf("Expected parameters: %v", expectedParams)
		t.Errorf("Parameter length failure: %d != %d", len(params), len(expectedParams))
		return
	}

	for i, v := range expectedParams {
		if v != params[i] {
			t.Errorf("Parsed   parameters: %v", params)
			t.Errorf("Expected parameters: %v", expectedParams)
			t.Errorf("Parameter parse failure: %s != %s at position %d", v, params[i], i)
		}
	}
}

func cursorSingleParamHelper(t *testing.T, command byte, funcName string) {
	t.Helper()
	funcCallParamHelper(t, []byte{command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'2', '3', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([23])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2])", funcName)})
}

func cursorTwoParamHelper(t *testing.T, command byte, funcName string) {
	t.Helper()
	funcCallParamHelper(t, []byte{command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1 1])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', '3', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([23 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2 3])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2 3])", funcName)})
}

func eraseHelper(t *testing.T, command byte, funcName string) {
	t.Helper()
	funcCallParamHelper(t, []byte{command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'1', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'3', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([3])", funcName)})
	funcCallParamHelper(t, []byte{'4', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'1', ';', '2', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
}

func scrollHelper(t *testing.T, command byte, funcName string) {
	t.Helper()
	funcCallParamHelper(t, []byte{command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'1', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'5', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([5])", funcName)})
	funcCallParamHelper(t, []byte{'4', ';', '6', command}, "CsiEntry", "Ground", []string{fmt.Sprintf("%s([4])", funcName)})
}

func clearOnStateChangeHelper(t *testing.T, start string, end string, bytes []byte) {
	t.Helper()
	p, _ := createTestParser(start)
	fillContext(p.context)
	if _, err := p.Parse(bytes); err != nil {
		t.Errorf("Error parsing %q: %v", string(bytes), err)
	}
	validateState(t, p.currState, end)
	validateEmptyContext(t, p.context)
}

func c0Helper(t *testing.T, bytes []byte, expectedState string, expectedCalls []string) {
	t.Helper()
	parser, evtHandler := createTestParser("Ground")
	if _, err := parser.Parse(bytes); err != nil {
		t.Errorf("Error parsing %q: %v", string(bytes), err)
	}
	validateState(t, parser.currState, expectedState)
	validateFuncCalls(t, evtHandler.FunctionCalls, expectedCalls)
}
