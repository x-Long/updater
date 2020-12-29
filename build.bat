@echo off

set work_dir=%CD%
if "%has_gopath%" == "" (
    set GOPATH=%GOPATH%;%cd%
    set has_gopath=1
)
pushd src
go build -ldflags="-w -s" updater.go
move updater.exe %work_dir%\ 
popd
upx\upx.exe -9 updater.exe
