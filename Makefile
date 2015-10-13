all:
	cd db/ && make
	go get ./...

test:
	cd pull/ && go test -v
