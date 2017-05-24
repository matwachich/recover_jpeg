@echo off
set GOOS=windows

set GOARCH=386
go build -o ImageSaver_x86.exe -ldflags "-w -s"

set GOARCH=amd64
go build -o ImageSaver_x64.exe -ldflags "-w -s"

pause