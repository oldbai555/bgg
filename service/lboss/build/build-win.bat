cd ..
go build -ldflags "-s -w"
upx -9 lboss.exe
