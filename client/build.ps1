$env:CGO_ENABLED="1"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags="-H windowsgui" -o client.exe  client.go

#linux go build -o server server.go
#$env:CGO_ENABLED="0"
#$env:GOOS="linux"
#$env:GOARCH="amd64"
#go build -o main main.go
 #>
