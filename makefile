lint:
	golangci-lint run -E revive,gocritic,gocyclo,goconst ./...

scan:
	./deslop scan . > results.txt 