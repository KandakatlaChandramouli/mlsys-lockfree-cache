#!/usr/bin/env bash

set -e

mkdir -p models

cd models

wget -O model.onnx \
https://huggingface.co/Xenova/all-MiniLM-L6-v2/resolve/main/onnx/model.onnx

wget -O onnxruntime-linux-x64-1.22.0.tgz \
https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz

tar -xzf onnxruntime-linux-x64-1.22.0.tgz

cp onnxruntime-linux-x64-1.22.0/lib/libonnxruntime.so* .

echo "model + runtime downloaded"
