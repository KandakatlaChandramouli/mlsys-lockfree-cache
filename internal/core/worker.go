
package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"fluxruntime/internal/model"
)

const (
	DefaultMaxBatchSize = 128
	DefaultBatchTimeout = 2 * time.Millisecond
	DefaultQueueDepth   = 4096

	NumShards = 8
)

type WorkerConfig struct {
	ShardID      int
	MaxBatchSize int
	BatchTimeout time.Duration
	ModelPath    string
}

type Worker struct {
	cfg     WorkerConfig
	model   *model.Model
	queue   chan *model.EmbeddingRequest
	running atomic.Bool
	wg      sync.WaitGroup
}

func NewWorker(
	cfg WorkerConfig,
) *Worker {

	if cfg.MaxBatchSize <= 0 {
		cfg.MaxBatchSize = DefaultMaxBatchSize
	}

	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = DefaultBatchTimeout
	}

	return &Worker{
		cfg: cfg,
		model: model.New(
			cfg.ModelPath,
		),
		queue: make(
			chan *model.EmbeddingRequest,
			DefaultQueueDepth,
		),
	}
}

func (w *Worker) Start(
	ctx context.Context,
) error {

	err := w.model.Load()

	if err != nil {
		return err
	}

	err = w.model.Warmup()

	if err != nil {
		return err
	}

	w.running.Store(true)

	w.wg.Add(1)

	go w.loop(ctx)

	return nil
}

func (w *Worker) Stop() {

	w.running.Store(false)

	close(w.queue)

	w.wg.Wait()

	w.model.Close()
}

func (w *Worker) Submit(
	req *model.EmbeddingRequest,
) error {

	select {

	case w.queue <- req:
		return nil

	default:
		return fmt.Errorf(
			"worker queue full",
		)
	}
}

func (w *Worker) loop(
	ctx context.Context,
) {

	defer w.wg.Done()

	batch := make(
		[]*model.EmbeddingRequest,
		0,
		w.cfg.MaxBatchSize,
	)

	for {

		batch = batch[:0]

		select {

		case <-ctx.Done():
			return

		case req, ok := <-w.queue:

			if !ok {
				return
			}

			batch = append(
				batch,
				req,
			)
		}

		timer := time.NewTimer(
			w.cfg.BatchTimeout,
		)

	drain:

		for len(batch) < w.cfg.MaxBatchSize {

			select {

			case req, ok := <-w.queue:

				if !ok {
					break drain
				}

				batch = append(
					batch,
					req,
				)

			case <-timer.C:
				break drain

			case <-ctx.Done():
				timer.Stop()
				return
			}
		}

		timer.Stop()

		err := w.model.RunBatch(
			&model.Batch{
				Requests: batch,
				Size:     len(batch),
			},
		)

		if err != nil {

			for _, req := range batch {

				req.ResultCh <- model.EmbeddingResult{
					Err: err,
				}
			}
		}
	}
}

type ShardedPool struct {
	workers []*Worker
	counter atomic.Uint64
}

func NewShardedPool(
	ctx context.Context,
	modelPath string,
) (*ShardedPool, error) {

	pool := &ShardedPool{
		workers: make(
			[]*Worker,
			NumShards,
		),
	}

	for i := 0; i < NumShards; i++ {

		w := NewWorker(
			WorkerConfig{
				ShardID:   i,
				ModelPath: modelPath,
			},
		)

		err := w.Start(ctx)

		if err != nil {
			return nil, err
		}

		pool.workers[i] = w
	}

	return pool, nil
}

func (p *ShardedPool) Embed(
	tokens []int32,
) ([]float32, error) {

	req := &model.EmbeddingRequest{
		Tokens: tokens,
		ResultCh: make(
			chan model.EmbeddingResult,
			1,
		),
	}

	idx := p.counter.Add(1)

	worker := p.workers[idx%uint64(len(p.workers))]

	err := worker.Submit(
		req,
	)

	if err != nil {
		return nil, err
	}

	result := <-req.ResultCh

	return result.Embedding, result.Err
}

func (p *ShardedPool) Shutdown() {

	for _, w := range p.workers {
		w.Stop()
	}
}
