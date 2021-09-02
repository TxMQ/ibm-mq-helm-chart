#!/bin/bash

# build chart

mkdir -p output

./chart-template.sh output
cp values.yaml output
cp -r templates output
