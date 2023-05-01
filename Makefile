ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
else
    SET=export
endif

.PHONY: all test

all :
	$(foreach I,$(wildcard coloring examples nternal/* keys simplehistory test/* tty*),pushd "$(I)" && go fmt && popd && ) go fmt
	go build

try :
	$(MAKE) all && "test/unicodetest/unicodetest"

test :
	go test -v

get :
	go get -u all
	go mod tidy

