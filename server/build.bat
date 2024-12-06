@echo off
setlocal

set OUTPUT=server.exe

go mod vendor
go build -o out/%OUTPUT% src/main.go
out\server.exe
