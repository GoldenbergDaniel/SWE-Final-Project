#!/bin/bash

OUTPUT="TradeEx"

mkdir -p out
go build -o out/$OUTPUT src/main.go
