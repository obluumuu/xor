.PHONY: fix-imports
fix-imports:
	gci write --skip-generated -s standard -s default -s localmodule .

.PHONY: lint-go
lint-go:
	golangci-lint run -v --config .golangci.yml

.PHONY: lint-proto
lint-proto:
	buf lint

.PHONY: lint
lint: fix-imports lint-proto lint-go


.PHONY: migrate-up
migrate-up:
	goose -dir migrations postgres "host=localhost user=app password=app dbname=app sslmode=disable" up


.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "host=localhost user=app password=app dbname=app sslmode=disable" down


.PHONY: gen-proto
gen-proto:
	protoc \
		--proto_path=proto \
		--go_out=gen/proto \
		--go_opt=Mproxy/proxy.proto=proxy/ \
		proto/**/**.proto

.PHONY: gen-sql
gen-sql:
	rm -rf gen/models && sqlc generate

.PHONY: gen
gen: gen-proto gen-sql


.PHONY: list-tests
list-tests:
	go test ./tests/... -list=.


.PHONY: test-production
test-production:
	go clean -testcache && TESTING_SERVER_TYPE=PRODUCTION go test -v ./tests/...

.PHONY: test-standalone
test-standalone:
	go clean -testcache && TESTING_SERVER_TYPE=STANDALONE go test -race -v ./tests/...

.PHONY: test
test: test-production test-standalone


.PHONY: install-tools
install-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.34.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install github.com/daixiang0/gci@v0.13.4
	go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.19.1
