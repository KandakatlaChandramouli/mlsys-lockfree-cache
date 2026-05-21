package ingest

import (
    "sync"
)

const BatchSize = 256

type Job struct {
    ID        string
    Embedding []float32
}

type Handler interface {
    InsertBatch(
        jobs []Job,
    )
}

type Pipeline struct {
    jobs chan Job
    wg   sync.WaitGroup
}

func New(
    workers int,
    queue int,
    h Handler,
) *Pipeline {

    p := &Pipeline{
        jobs: make(
            chan Job,
            queue,
        ),
    }

    for i := 0; i < workers; i++ {

        p.wg.Add(1)

        go func() {

            defer p.wg.Done()

            batch := make(
                []Job,
                0,
                BatchSize,
            )

            flush := func() {

                if len(batch) == 0 {
                    return
                }

                h.InsertBatch(
                    batch,
                )

                batch = batch[:0]
            }

            for job := range p.jobs {

                batch = append(
                    batch,
                    job,
                )

                if len(batch) >= BatchSize {
                    flush()
                }
            }

            flush()
        }()
    }

    return p
}

func (p *Pipeline) Submit(
    job Job,
) {

    p.jobs <- job
}

func (p *Pipeline) Close() {

    close(p.jobs)

    p.wg.Wait()
}
