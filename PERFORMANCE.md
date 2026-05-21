# FluxRuntime Performance Evolution

## Scalar Go
~56ms

## Scalar Assembly
~57ms

## AVX2 YMM SIMD
~8.5ms

## AVX2 FMA
~5-6ms expected

## Optimization Stack
- YMM vector lanes
- fused multiply-add
- horizontal reduction
- zero allocation traversal
- ring-buffer frontier
- SIMD cosine kernels
