#!/bin/bash

OUTPUT="TradEx"

mkdir -p out
go mod vendor
go build -o out/$OUTPUT src/main.go
out/TradEx
