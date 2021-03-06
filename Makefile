all: lint build test

build:
	go build ./...

install:
	go get github.com/alecthomas/gometalinter && gometalinter --install --update

lint:
	gometalinter --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

cover:
	go test -coverprofile=coverage.out . && go tool cover -html=coverage.out && rm -f coverage.out

test:
	go test -cover .
