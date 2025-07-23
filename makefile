.PHONY: lint
lint:
	staticcheck ./...
	

.PHONY: test
PKGS := $(shell go list ./... | grep -vE "/test")

.PHONY: test
test:
	go test $(PKGS)

.PHONY: coverage
coverage:
	go test $(PKGS) -coverprofile=coverage.out

.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html

.PHONY: coverage-total
coverage-total: coverage
	go tool cover -func=coverage.out | grep total | awk '{print $$3}'