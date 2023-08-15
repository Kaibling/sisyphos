build:
	go mod tidy && go build -o sisyphos -buildvcs=false
test:
	go test -v ./... -race -covermode=atomic
run: deps
	air server

deps:
	go install github.com/cosmtrek/air@latest

run-ui:
	cd ui;	npm run dev -- --host 0.0.0.0

start-docker:
	docker compose down
	docker compose up --build --remove-orphans