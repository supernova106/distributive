TEST?=$(glide novendor)
NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)

all: deps build

deps:
	glide install

updatedeps:
	glide update

build: deps
	@mkdir -p bin/
	go build -o bin/$(NAME)

test: deps
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

package: build
	tar -zcvf bin/$(NAME).tar.gz bin/$(NAME)

.PHONY: all deps updatedeps build test package