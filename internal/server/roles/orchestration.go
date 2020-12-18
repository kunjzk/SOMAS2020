package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

// callVote possible implementation of voting
func callVote(speakerID int, whateverIsBeingVotedOn string) {
	// Do voting

	noIslandAlive := rules.VariableValuePair{
		VariableName: "no_islands_alive",
		Values:       []float64{5},
	}
	noIslandsVoting := rules.VariableValuePair{
		VariableName: "no_islands_voted",
		Values:       []float64{5},
	}
	err := updateTurnHistory(speakerID, []rules.VariableValuePair{noIslandAlive, noIslandsVoting})
	if err != nil {
		// exit with error
	} else {
		// carry on
	}
}
