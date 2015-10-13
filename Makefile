all:
	cd db/ && make
	cd pull/ && go build

test:
	cd pull/ && go test -v
