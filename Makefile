ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
else
    SHELL=bash
    SET=export
endif

.PHONY: all test

all :
	$(foreach I,$(wildcard coloring examples internal/* keys simplehistory test/* tty*),pushd "$(I)" && go fmt && popd && ) go fmt
	go build

demo :
	go run test/unicodetest/main.go

test :
	go test -v

get :
	go get -u all
	go mod tidy

