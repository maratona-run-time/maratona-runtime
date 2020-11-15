export CGO_ENABLED=0
diff -u <(echo -n) <(gofmt -d ./)
go test ./... -cover -v
