all:
	$(MAKE) -C db
	$(MAKE) -C oauth
	go get ./...

test:
	go get github.com/jarcoal/httpmock
	go get github.com/go-sql-driver/mysql
	cd pull/ && go test -v
