.PHONY: lint, test, fmtall, all, run, run_sender, run_listener

lint:
	cd listener && golangci-lint run
	cd sender && golangci-lint run

test:
	cd listener && go test -race
	cd sender && go test -race

fmtall:
	cd listener && go fmt ./...
	cd sender && go fmt ./...

all:
	make fmtall
	make lint
	make test

run:
	make run_listener
	make run_sender

run_sender:
	cd sender && go build -o send
	mv sender/send ./
	./send -file=sample10mb.txt

run_listener:
	cd listener && go build -o listen
	mv listener/listen ./
	./listen

