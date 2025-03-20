go env -w CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o uploader ./
unset GOOS
unset GOARCH
go build -o uploader.exe ./
go env -w CGO_ENABLED=1
