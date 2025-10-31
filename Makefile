ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    NUL=NUL
    WHICH=where.exe
else
    SHELL=bash
    SET=export
    NUL=/dev/null
    WHICH=which
endif

ifndef GO
    SUPPORTGO=go1.20.14
    GO:=$(shell $(WHICH) $(SUPPORTGO) 2>$(NUL)|| echo go)
endif

.PHONY: all test bench

all :
	$(GO) fmt ./...
	$(GO) build

demo :
	$(GO) run test/unicodetest/main.go

demo-future :
	$(GO) run -tags=tty10,orgxwidth test/unicodetest/main.go

test :
	$(GO) test -v ./...

get :
	$(GO) get -u all
	$(GO) mod tidy

$(SUPPORTGO):
	go install golang.org/dl/$(SUPPORTGO)@latest
	$(SUPPORTGO) download

bench :
	cd test/bench && lispect time.lsp
