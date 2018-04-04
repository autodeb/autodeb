all: fmt list-packages-to-update vet lint

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile

GO_IMPORT_PATH := salsa.debian.org/aviau/autopkgupdate

list-packages-to-update: $(SOURCES)
	go build -v -o list-packages-with-newer-upstream-versions ${GO_IMPORT_PATH}/cmd/list-packages-with-newer-upstream-versions

.PHONY: vet
vet: install
	go vet -v ${GO_IMPORT_PATH}/...

.PHONY: install
install:
	go install ${GO_IMPORT_PATH}/...

.PHONY: lint
lint:
	golint ${GO_IMPORT_PATH}/...

.PHONY: fmt
fmt:
	go fmt ${GO_IMPORT_PATH}/...


.PHONY: get-deps
get-deps:
	go get -t ${GO_IMPORT_PATH}/...


.PHONY: clean
clean:
	rm -f list-packages-with-newer-upstream-versions
