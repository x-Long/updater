set GOPATH=%GOPATH%;%cd%
cd src\updater
go build -ldflags="-w -s" updater.go
.\upx\upx.exe -9 updater.exe