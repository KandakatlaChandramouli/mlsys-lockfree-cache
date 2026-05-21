# FluxRuntime Release Notes

## Runtime
- ONNX Runtime integration
- Dynamic batching
- Worker sharding
- TCP transport optimization
- Embedding pooling
- Request metrics

## Retrieval
- Vector store
- Graph ANN index
- HNSW traversal
- RAG engine
- Persistence layer

## Optimization
- SIMD cosine similarity
- Zero-allocation search
- Ring-buffer traversal
- Candidate pruning
- TopK optimization

## Benchmarks

Inference:
- ~22k req/sec
- p95 ~12ms
- p99 ~16ms

ANN:
- ~9ms search latency
- 100k vectors

## Status
Production-grade experimental MLSys runtime prototype.
