package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	tmcoreypes "github.com/tendermint/tendermint/rpc/coretypes"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"

	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

const (
	keyHeight = "height"
	keyModule = "module"
)

func (w *Worker) process(ctx context.Context, workerIndex int, recoverMode bool) {
	var parsedCount int
	defer w.wg.Done()
	defer func() {
		w.log.Debug().Msgf("worker: %d. parsed %d blocks", workerIndex, parsedCount)
	}()

	for height := range w.heightCh {
		select {
		case <-ctx.Done():
			w.log.Info().Int("worker_index", workerIndex).Msg("done worker")
			return
		default:
		}

		// for debug
		parsedCount++

		w.processHeight(ctx, workerIndex, height, recoverMode)
	}
}

func (w *Worker) processHeight(ctx context.Context, workerIndex int, height int64, recoveryMode bool) { // nolint:gocognit
	if recoveryMode {
		defer func() {
			if r := recover(); r != nil {
				w.setErrorStatusWithLogging(ctx, height, fmt.Sprint(r))
				w.log.Error().Int64("height", height).Msgf("panic occurred!\n%v", r)
			}
		}()
	}

	if err := w.checkOrCreateBlockInStorage(ctx, height); err != nil {
		switch {
		case errors.Is(err, ErrBlockProcessed):
			w.log.Debug().Int64(keyHeight, height).Msg("block already processed. skip height")
		case errors.Is(err, ErrBlockProcessing):
			w.log.Debug().Int64(keyHeight, height).Msg("block is already processing now. skip height")
		case errors.Is(err, ErrBlockError):
			w.log.Debug().Int64(keyHeight, height).Msg("block processed with error. " +
				"if you want to process this height again see PROCESS_ERROR_BLOCKS ENV")
		}

		return
	}

	if height == 0 {
		w.log.Info().Int("worker_number", workerIndex).Msg("Parse genesis")

		_genesisDur := time.Now()

		genesis, err := w.rpcClient.Genesis(ctx)
		if err != nil {
			w.setErrorStatusWithLogging(ctx, height, err.Error())
			w.log.Error().Err(err).Msg("get genesis error")
			return
		}

		w.log.Debug().Int("worker_number", workerIndex).
			Dur("duration", time.Since(_genesisDur)).
			Msg("get genesis")

		if err = w.processGenesis(ctx, genesis); err != nil {
			w.log.Error().Err(err).Msg("processHeight genesis error")
			w.setErrorStatusWithLogging(ctx, height, err.Error())
			return
		}

		if err = w.storage.SetProcessedStatus(ctx, height); err != nil {
			w.log.Error().Err(err).Int64(keyHeight, height).Msg("can't set processed status in storage")
		}

		return
	}

	w.log.Info().Int("worker_number", workerIndex).Int64("height", height).Msg("parse block")

	g, ctx2 := errgroup.WithContext(ctx)

	var block *tmcoreypes.ResultBlock

	g.Go(func() error {
		var err error

		_blockDur := time.Now()
		if block, err = w.grpcClient.Block(ctx2, height); err != nil {
			return fmt.Errorf("failed to get block: %w", err)
		}
		w.log.Debug().
			Int("worker_number", workerIndex).
			Int64("block_height", height).
			Dur("get_block_dur", time.Since(_blockDur)).
			Msg("get block info")
		return nil
	})

	if err := g.Wait(); err != nil {
		w.log.Error().Int64(keyHeight, height).Err(err).Msg("processHeight error")
		w.setErrorStatusWithLogging(ctx, height, err.Error())
		return
	}

	_txsDur := time.Now()

	txsRes, err := w.grpcClient.Txs(ctx, height, block.Block.Data.Txs)
	if err != nil {
		w.log.Error().Err(err).Msg("get txs error")
		w.setErrorStatusWithLogging(ctx, height, err.Error())
		return
	}

	w.log.Debug().
		Int("worker_number", workerIndex).
		Int64("block_height", height).
		Dur("txs_dur", time.Since(_txsDur)).
		Msg("Get txs info")

	txs := types.NewTxsFromTmTxs(txsRes, w.cdc)
	g, ctx2 = errgroup.WithContext(ctx)

	g.Go(func() error {
		return w.withMetrics("block", func() error {
			return w.processBlock(ctx2, types.NewBlockFromTmBlock(block, txs.TotalGas()))
		})
	})
	g.Go(func() error {
		return w.withMetrics("txs", func() error {
			return w.processTxs(ctx2, txs)
		})
	})

	if err := g.Wait(); err != nil {
		w.setErrorStatusWithLogging(ctx, height, err.Error())
		return
	}

	if err := w.storage.SetProcessedStatus(ctx, height); err != nil {
		w.log.Error().Err(err).Int64(keyHeight, height).Msg("can't set processed status in storage")
	}
}

func (w *Worker) processGenesis(ctx context.Context, genesis *tmtypes.GenesisDoc) error {
	var appState map[string]json.RawMessage
	if err := jsoniter.Unmarshal(genesis.AppState, &appState); err != nil {
		w.log.Err(err).Msg("error unmarshalling genesis doc")
		return err
	}

	for _, m := range genesisHandlers {
		if err := m.HandleGenesis(ctx, genesis, appState); err != nil {
			w.log.Error().Err(err).Str(keyModule, m.Name()).Msg("handle genesis error")
		}
	}

	return nil
}

func (w *Worker) processBlock(ctx context.Context, block *types.Block) error {
	for _, m := range blockHandlers {
		if err := m.HandleBlock(ctx, block); err != nil {
			w.log.Error().
				Err(err).
				Int64(keyHeight, block.Height).
				Str(keyModule, m.Name()).
				Msg("HandleBlock error")

			return err
		}
	}

	return nil
}

func (w *Worker) processTxs(ctx context.Context, txs []*types.Tx) error {
	for _, tx := range txs {
		for _, m := range transactionHandlers {
			if err := m.HandleTx(ctx, tx); err != nil {
				w.log.Error().
					Err(err).
					Int64(keyHeight, tx.Height).
					Str(keyModule, m.Name()).
					Msg("HandleTX error")

				return err
			}
		}
	}

	return nil
}
