default: vet test

vet:
	go vet ./...

test:
	go test ./...

bench:
	go test ./... -run=NONE -bench=. -benchmem
