GIT_HASH:=$(shell git rev-parse --short HEAD)
DIRTY:=$(shell test -z "`git status --porcelain`" || echo "-dirty")
VERSION:=$(GIT_HASH)$(DIRTY)
TIME:=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BIN:=huectl
PKG:=github.com/callebjorkell/huectl

.PHONY: bin install

bin:
	go build -ldflags "-X $(PKG)/cmd.commit=$(VERSION) -X $(PKG)/cmd.time=$(TIME)" -o $(BIN) .

install: bin
	cp $(BIN) $(HOME)/bin/

clean:
	rm -f $(BIN)