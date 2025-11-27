COVERAGE_FILE := coverage.txt
COVERAGE_HTML := coverage.html

benchmark:
	go test -bench=. -benchmem

test:
	go test -v -count=1 -race . -cover

test\:coverage:
	go test -v -count=1 -race . -coverprofile=$(COVERAGE_FILE)


coverage-html:
	@if ! go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML); then \
		echo "Error: no coverage file, ensure you run 'make test:coverage' first"; \
		exit 1; \
	fi

coverage-clean:
	@echo "Coverage files cleaned"; \
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

lint:	
	golangci-lint run



