ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
else
    SHELL=bash
    SET=export
endif

.PHONY: all test

all :
	go fmt ./...
	go build

demo :
	go run test/unicodetest/main.go

demo-future :
	go run -tags=tty10,orgxwidth test/unicodetest/main.go

test :
	pushd "internal/moji" && go test -v && popd
	go test -v

get :
	go get -u all
	go mod tidy

