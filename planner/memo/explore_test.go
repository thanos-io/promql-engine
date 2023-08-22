package memo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var explorationMarkTestCases = []struct {
	input    ExplorationMark // The initial exploration mark.
	method   string
	round    ExplorationRound
	expected ExplorationMark // The expected exploration mark.
}{
	{
		input:    ExplorationMark(0),
		method:   "set",
		round:    ExplorationRound(0),
		expected: ExplorationMark(0b1),
	},
	{
		input:    ExplorationMark(0),
		method:   "set",
		round:    ExplorationRound(1),
		expected: ExplorationMark(0b10),
	},
	{
		input:    ExplorationMark(0),
		method:   "set",
		round:    ExplorationRound(2),
		expected: ExplorationMark(0b100),
	},
	{
		input:    ExplorationMark(0b1110),
		method:   "unset",
		round:    ExplorationRound(2),
		expected: ExplorationMark(0b1010),
	},
	// TODO add tests
}

func TestExplorationMark(t *testing.T) {
	for _, test := range explorationMarkTestCases {
		description := fmt.Sprintf("exploration mark %b after %s on round %d should be %b", test.input, test.method, test.round, test.expected)
		t.Run(description, func(t *testing.T) {
			mark := test.input
			switch test.method {
			case "set":
				mark.SetExplore(test.round, true)
			case "unset":
				mark.SetExplore(test.round, false)
			default:
				require.Fail(t, "unrecognized method %s", test.method)
			}
			require.Equal(t, test.expected, mark, "error on input '%s'", test.input)
		})
	}
}
