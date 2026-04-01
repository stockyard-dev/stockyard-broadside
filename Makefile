build:
	CGO_ENABLED=0 go build -o broadside ./cmd/broadside/

run: build
	./broadside

test:
	go test ./...

clean:
	rm -f broadside

.PHONY: build run test clean
