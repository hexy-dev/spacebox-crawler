package worker

import (
	"context"

	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

var (
	ErrBlockProcessed  = errors.New("block already processed")
	ErrBlockProcessing = errors.New("block is processing right now")
	ErrBlockError      = errors.New("block processed with error")
)

func (w *Worker) setErrorStatusWithLogging(ctx context.Context, height int64, msg string) {
	if err := w.storage.SetErrorStatus(ctx, height, msg); err != nil {
		w.log.Error().Err(err).Int64("height", height).Msg("can't set error status in storage")
	}
}

func (w *Worker) checkOrCreateBlockInStorage(ctx context.Context, height int64) error {
	block, err := w.storage.GetBlockByHeight(ctx, height)
	if err != nil && errors.Is(err, types.ErrBlockNotFound) {
		// create new block
		if err = w.storage.CreateBlock(ctx, w.tsM.NewBlock(height)); err != nil {
			w.log.Error().Err(err).Int64("height", height).Msg("can't create new block in storage")
			return err
		}

		return nil
	} else if err != nil {
		// got some error from storage
		return err
	}

	// block exists check status
	switch {
	// block info already in kafka
	case block.Status.IsProcessed():
		return ErrBlockProcessed
	// block now is processing
	case block.Status.IsProcessing():
		return ErrBlockProcessing
	// block processed with error, skip if needed
	case block.Status.IsError():
		if !w.cfg.ProcessErrorBlocks {
			return ErrBlockError
		}
		return w.storage.UpdateStatus(ctx, height, model.StatusProcessing)
	}
	return nil
}
