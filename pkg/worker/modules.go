package worker

import "github.com/bro-n-bro/spacebox-crawler/v2/types"

var (
	transactionHandlers []types.TransactionHandler
	blockHandlers       []types.BlockHandler
	genesisHandlers     []types.GenesisHandler
)

// fillModules fills the module handlers.
func (w *Worker) fillModules() {
	for _, module := range w.modules {
		if tI, ok := module.(types.TransactionHandler); ok {
			transactionHandlers = append(transactionHandlers, tI)
		}
		if bI, ok := module.(types.BlockHandler); ok {
			blockHandlers = append(blockHandlers, bI)
		}
		if gI, ok := module.(types.GenesisHandler); ok {
			genesisHandlers = append(genesisHandlers, gI)
		}
	}
}
