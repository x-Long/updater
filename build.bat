@echo off

set work_dir=%CD%

REM for running on windows x86 os we need compile with x86 go.exe
if "%go32_set%" == "" (
    set go32_path=D:\GitRepos\longway\go
    set go32_set=1
)

if "%gopath_set%" == "" (
    set GOPATH=%cd%;%GOPATH%
    set gopath_set=1
)
pushd src
if exist %go32_path%\bin\go.exe (
    set go32_bin=%go32_path%"\bin\go.exe"
) else (
    set go32_bin=go.exe
)
echo go path is %go32_bin%
%go32_bin% build -ldflags="-w -s" updater.go
move updater.exe %work_dir%\ 
popd
upx\upx.exe -9 updater.exe

