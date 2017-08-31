PKG=$(shell go list | grep -v 'vendor')

default: vet test

vet:
	go vet $(PKG)

test:
	go test $(PKG)

bench:
	go test $(PKG) -run=NONE -bench=. -benchmem
