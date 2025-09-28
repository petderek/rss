.PHONY: all test clean
all: refuss

test:
	go test ./...

clean:
	rm refuss

refuss: *.go ./cmd/refuss/*
	go build -o refuss ./cmd/refuss


