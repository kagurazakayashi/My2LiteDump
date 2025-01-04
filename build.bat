SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
DEL My2LiteDump.xz
go build .
xz -z -e -9 -T 0 -v My2LiteDump
