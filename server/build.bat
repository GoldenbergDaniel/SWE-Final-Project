@echo off
setlocal

go build -o out/server.exe src2/main.go
out\server.exe