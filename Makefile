ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
else
    SET=export
endif

.PHONY: all test

all :
	go fmt
	go build
	cd "test/unicodetest" && go fmt && go build

test :
	$(MAKE) all && "test/unicodetest/unicodetest"

get :
	go get -u all
	go mod tidy

