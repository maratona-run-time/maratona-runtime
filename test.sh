set -e
export CGO_ENABLED=0
diff -u <(echo -n) <(gofmt -d ./)

test_folder() {
    go test $1 -v -cover | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
    go test $1 &> /dev/null
}

test_folder github.com/maratona-run-time/Maratona-Runtime/verdict/src
test_folder github.com/maratona-run-time/Maratona-Runtime/compiler/src
test_folder github.com/maratona-run-time/Maratona-Runtime/executor/src
test_folder github.com/maratona-run-time/Maratona-Runtime/comparator/src
