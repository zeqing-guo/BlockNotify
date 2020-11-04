unamestr := $(shell uname)

ifeq ($(unamestr),Linux)
	platform := linux
else
  ifeq ($(unamestr),Darwin)
	platform := darwin
  endif
endif

tags := $(platform)

GOBUILD := go build -tags "$(tags)"

bin/blocknotify:
	$(GOBUILD) \
		-o bin/blocknotify \
		github.com/zeqing-guo/BlockNotify/cmd/blocknotify

all: bin/blocknotify

clean:
	rm -rf bin/*