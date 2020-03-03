BINNAME=rgnes

.PHONY: all clean

all: test build

build:
	go build -ldflags='-w -s' -o $(BINNAME) ./cmd/$(BINNAME)/main.go

test:
	go test -race -v ./...

clean:
	go clean
	go clean -testcache
	rm -f $(BINNAME)