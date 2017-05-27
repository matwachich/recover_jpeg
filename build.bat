@echo off
set GOOS=windows

set GOARCH=386
go build -o ImageSaver_x86.exe -ldflags "-w -s"
upx ImageSaver_x86.exe

set GOARCH=amd64
go build -o ImageSaver_x64.exe -ldflags "-w -s"
upx ImageSaver_x64.exe

pause