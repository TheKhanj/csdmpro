GO_FILES = $(shell find core -type f -name '*.go') \
						$(shell find tg -type f -name '*.go') \
						$(shell find db -type f -name '*.go') \
						main.go wire.go

all: build

wire: $(GO_FILES)
	wire

build: $(GO_FILES) wire
	go build

build_watch: $(GO_FILES)
	printf "%s\n" $(GO_FILES) | \
    entr -c -c $(MAKE) build

run: build
	./csdmpro

run_watch: $(GO_FILES)
	printf "%s\n" $(GO_FILES) | \
    entr -s -r -c -c "$(MAKE) run"
