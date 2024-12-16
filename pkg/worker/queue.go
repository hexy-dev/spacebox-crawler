package worker

import (
	"context"
	"sync"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

func (w *Worker) enqueueHeight(ctx context.Context, wg *sync.WaitGroup, startHeight, stopHeight int64) {
	defer wg.Done()

	w.log.Debug().Msgf("try to parse: %d count of blocks", stopHeight-startHeight)

	ctx, w.stopEnqueueHeight = context.WithCancel(ctx)
	defer w.stopEnqueueHeight()

	// send zero height to process genesis
	if startHeight != 0 && w.cfg.ProcessGenesis {
		w.heightCh <- 0
	}

	for height := startHeight; height >= 0 && height <= stopHeight; height++ {
		// safe from closed channel
		select {
		case <-ctx.Done():
			w.log.Info().Msg("stop enqueueHeight")
			return
		case w.heightCh <- height: // put height to channel for processing the block
		}
	}
}

func (w *Worker) enqueueNewBlocks(ctx context.Context, eventCh <-chan coretypes.ResultEvent) {
	ctx, w.stopWsListener = context.WithCancel(ctx)
	defer w.stopWsListener()
	w.log.Info().Msg("listening for new block events")

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("stop new block listener")
			return
		case e := <-eventCh:
			newBlock, ok := e.Data.(types.EventDataNewBlock)
			if !ok {
				w.log.Warn().Msg("failed to cast ws event to EventDataNewBlock type")
				continue
			}

			height := newBlock.Block.Header.Height
			w.log.Info().Int64("height", height).Msg("enqueueing new block")

			select {
			case <-ctx.Done():
				w.log.Info().Msg("stop new block listener")
				return
			case w.heightCh <- height:
			}
		}
	}
}

func (w *Worker) enqueueErrorBlocks(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(w.cfg.ProcessErrorBlocksInterval)
	defer ticker.Stop()

	ctx, w.stopEnqueueErrorBlocks = context.WithCancel(ctx)
	defer w.stopEnqueueErrorBlocks()

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("stop GetErrorBlockHeights")
			return
		case <-ticker.C:
			heights, err := w.storage.GetErrorBlockHeights(ctx)
			if err != nil {
				w.log.Error().Err(err).Str("func", "GetErrorBlockHeights").Msg("can't enqueueErrorBlocks")
				return
			}

			for _, height := range heights {
				// safe from closed channel
				select {
				case <-ctx.Done():
					w.log.Info().Msg("stop GetErrorBlockHeights")
					return
				case w.heightCh <- height:
				}
			}
		}
	}
}
