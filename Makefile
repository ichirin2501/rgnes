BINNAME=rgnes

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: build
build:
	go build -ldflags='-w -s' -o $(BINNAME) ./cmd/$(BINNAME)/main.go

#.PHONY: test
#test:
#	go test -race -v ./...

.PHONY: loop
loop:
	@while true; do sleep 1; done

.PHONY: test
test:
	docker build -t rgnes -f Dockerfile.test .
	docker run --rm rgnes:latest go test -race -v ./...

.PHONY: clean
clean:
	go clean
	go clean -testcache
	rm -f $(BINNAME)
