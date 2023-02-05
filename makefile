build:
	go mod tidy && go build -o sisyphos
test:
	go test -v ./... -race -covermode=atomic
run:
	air server