GO_FILES = $(shell find core -type f -name '*.go') \
						$(shell find tg -type f -name '*.go') \
						$(shell find db -type f -name '*.go') \
						main.go wire.go

DEV_GO_FILES = $(shell [ -f .dev ] && find ../tgool -type f -name '*.go')

all: build

wire: $(GO_FILES) $(DEV_GO_FILES)
	wire

build: $(GO_FILES) $(DEV_GO_FILES) wire
	go build

build_watch: $(GO_FILES)
	printf "%s\n" $(GO_FILES) $(DEV_GO_FILES) | \
    entr -c -c $(MAKE) build

run: build
	./csdmpro

run_watch: $(GO_FILES)
	printf "%s\n" $(GO_FILES) | \
    entr -s -r -c -c "$(MAKE) run"

.dev:
		echo 'replace github.com/thekhanj/tgool => ../tgool' >> go.mod
		go mod tidy
		touch .dev

dev: .dev

undev:
	sed -i '/replace .*tgool/d' go.mod
	go mod tidy
	rm .dev

.PHONY: build build_watch run run_watch dev undev
