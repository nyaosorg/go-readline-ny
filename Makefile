ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    D=$\\
else
    SET=export
    D=/
endif

.PHONY: all test

all :
	go fmt
	go build
	cd cmd$(D)unicodetest && go fmt && go build

test :
	$(MAKE) all && cmd$(D)unicodetest$(D)unicodetest

get :
	go get -u all
	go mod tidy

