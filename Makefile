##
## This makefile is responsible for building, linting and testing the project
##

.PHONY: all
all: fmt \
	list-packages-with-newer-upstream-versions \
	update-random-package \
	vet \
	lint

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile

GO_IMPORT_PATH := salsa.debian.org/aviau/autopkgupdate

list-packages-with-newer-upstream-versions: $(SOURCES) install
	go build -v -o list-packages-with-newer-upstream-versions ${GO_IMPORT_PATH}/cmd/list-packages-with-newer-upstream-versions

update-random-package: $(SOURCES) install
	go build -v -o update-random-package ${GO_IMPORT_PATH}/cmd/update-random-package

.PHONY: vet
vet: install
	go vet -v ${GO_IMPORT_PATH}/...

.PHONY: install
install:
	go install ${GO_IMPORT_PATH}/...

.PHONY: lint
lint:
	${GOPATH}/bin/golint ${GO_IMPORT_PATH}/...

.PHONY: fmt
fmt:
	go fmt ${GO_IMPORT_PATH}/...


.PHONY: get-deps
get-deps:
	go get -t ${GO_IMPORT_PATH}/...

.PHONY: clean
clean:
	rm -f list-packages-with-newer-upstream-versions
	rm -f update-random-package
	rm -rf update-random-package-output-*
