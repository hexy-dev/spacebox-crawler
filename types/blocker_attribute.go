package types

import abci "github.com/cometbft/cometbft/abci/types"

type BlockerEvents map[string][]abci.Event

func NewBlockerEventsAttributes(events []abci.Event) BlockerEvents {
	res := make(BlockerEvents)
	for i := 0; i < len(events); i++ {
		res[events[i].Type] = append(res[events[i].Type], events[i])
	}
	return res
}
