diff -u <(echo -n) <(gofmt -d ./)
go test github.com/maratona-run-time/Maratona-Runtime/verdict/src -cover -v
go test github.com/maratona-run-time/Maratona-Runtime/executor/src -cover -v
go test github.com/maratona-run-time/Maratona-Runtime/compiler/src -cover -v
go test github.com/maratona-run-time/Maratona-Runtime/comparator/src -cover -v