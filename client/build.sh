!#/bin/sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myngrok

#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build *.go