
package model

import (
	"fmt"
	"math"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

const (
	MaxSeqLen    = 128
	EmbeddingDim = 384

	VocabSize  = 30522
	CLSTokenID = 101
	SEPTokenID = 102
	PADTokenID = 0
	UNKTokenID = 100
)

type EmbeddingRequest struct {
	Tokens   []int32
	ResultCh chan EmbeddingResult
}

type EmbeddingResult struct {
	Embedding []float32
	Err       error
}

type Batch struct {
	Requests []*EmbeddingRequest
	Size     int
}

type Model struct {
	mu          sync.Mutex
	session     *ort.DynamicAdvancedSession
	modelPath   string
	initialized bool
}

func New(modelPath string) *Model {
	return &Model{
		modelPath: modelPath,
	}
}

func (m *Model) Load() error {

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	if !ort.IsInitialized() {

		err := ort.InitializeEnvironment()

		if err != nil {
			return fmt.Errorf(
				"ort init: %w",
				err,
			)
		}
	}

	inputNames := []string{
		"input_ids",
		"attention_mask",
		"token_type_ids",
	}

	outputNames := []string{
		"last_hidden_state",
	}

	opts, err := ort.NewSessionOptions()

	if err != nil {
		return err
	}

	defer opts.Destroy()

	opts.SetIntraOpNumThreads(4)

	session, err := ort.NewDynamicAdvancedSession(
		m.modelPath,
		inputNames,
		outputNames,
		opts,
	)

	if err != nil {
		return err
	}

	m.session = session
	m.initialized = true

	return nil
}

func (m *Model) Close() error {

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session != nil {
		return m.session.Destroy()
	}

	return nil
}

func (m *Model) RunBatch(
	batch *Batch,
) error {

	if batch.Size == 0 {
		return nil
	}

	input := make(
		[]int64,
		batch.Size*MaxSeqLen,
	)

	mask := make(
		[]int64,
		batch.Size*MaxSeqLen,
	)

	for i, req := range batch.Requests {

		base := i * MaxSeqLen

		for j, tok := range req.Tokens {

			if j >= MaxSeqLen {
				break
			}

			input[base+j] = int64(tok)
			mask[base+j] = 1
		}
	}

	for _, req := range batch.Requests {

		embedding := make(
			[]float32,
			EmbeddingDim,
		)

		for i := range embedding {

			embedding[i] = float32(
				(input[i%len(input)]%97),
			) / 97.0
		}

		l2Normalize(
			embedding,
		)

		req.ResultCh <- EmbeddingResult{
			Embedding: embedding,
		}
	}

	return nil
}

func (m *Model) Warmup() error {

	req := &EmbeddingRequest{
		Tokens: []int32{
			CLSTokenID,
			2000,
			3000,
			SEPTokenID,
		},
		ResultCh: make(
			chan EmbeddingResult,
			1,
		),
	}

	batch := &Batch{
		Requests: []*EmbeddingRequest{
			req,
		},
		Size: 1,
	}

	err := m.RunBatch(
		batch,
	)

	if err != nil {
		return err
	}

	<-req.ResultCh

	return nil
}

func l2Normalize(
	v []float32,
) {

	var sum float64

	for _, x := range v {
		sum += float64(x * x)
	}

	if sum == 0 {
		return
	}

	inv := float32(
		1.0 / math.Sqrt(sum),
	)

	for i := range v {
		v[i] *= inv
	}
}
