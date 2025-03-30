GO_FILES = $(shell find . -type f -name '*.go') \
					 $(shell find ../tgool -type f -name '*.go')

all: build

wire: $(GO_FILES)
	wire

build: $(GO_FILES) wire
	go build

build_watch: $(GO_FILES)
	printf "%s\n" $(GO_FILES) | \
    entr -c -c $(MAKE) build
