default: test

test:
	go test ./...

bench:
	go test ./... -run=NONE -bench=. -benchmem

staticcheck:
	staticcheck ./...
