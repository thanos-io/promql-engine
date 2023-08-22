package memo

const (
	InitialExplorationRound ExplorationRound = 0
)

type ExplorationRound int // The exploration round.
type ExplorationMark int  // The exploration mark, if i-th bit of the mark is set, then i-th round is explored.

func (e *ExplorationMark) IsExplored(round ExplorationRound) bool {
	return (*e & (1 << round)) != 0
}

func (e *ExplorationMark) SetExplore(round ExplorationRound, explored bool) {
	if explored {
		*e |= 1 << round
	} else {
		*e &= ^(1 << round)
	}
}
