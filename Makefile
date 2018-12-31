LATEST_TAG := $(shell git describe --abbrev=0 --tags)

setup:
	go get \
		github.com/Songmu/goxz/cmd/goxz \
		github.com/tcnksm/ghr \
		github.com/golang/lint/golint \
		github.com/golang/dep/cmd/dep
	go get -d -t ./...
	dep ensure

test:
	go test -v ./...

lint:
	go vet ./...
	golint -set_exit_status ./...

dist:
	goxz -d dist/$(LATEST_TAG) -z -os windows,darwin,linux -arch amd64,386
	goxz -d dist/$(LATEST_TAG) -z -os linux -arch mipsle,arm

clean:
	rm -rf dist/*

release:
	ghr -u papix -r mackerel-plugin-nature-remo $(LATEST_TAG) dist/$(LATEST_TAG)

.PHONY: setup test lint dist clean
