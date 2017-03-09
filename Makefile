########################################################################################

.PHONY = fmt all clean deps deps-test test

########################################################################################

all: shdoc

shdoc:
	go build shdoc.go

deps:
	go get -v pkg.re/essentialkaos/ek.v7

deps-test:
	go get -v pkg.re/check.v1

test:
	go test -covermode=count .

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

clean:
	rm -f shdoc

########################################################################################

