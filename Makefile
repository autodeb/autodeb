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
	golint ${GO_IMPORT_PATH}/...

.PHONY: fmt
fmt:
	go fmt ${GO_IMPORT_PATH}/...


.PHONY: get-deps
get-deps:
	go get -t ${GO_IMPORT_PATH}/...

.PHONY: clean
clean: docker-clean
	rm -f list-packages-with-newer-upstream-versions
	rm -f update-random-package
	rm -rf update-random-package-output-*

####################
## Docker targets ##
####################

.PHONY: docker-clean
docker-clean:
	make -f Makefile_docker clean

.PHONY: .docker-run-makefile-target
.docker-golang-run-makefile-target:
	MAKEFILE_TARGET=${MAKEFILE_TARGET} make -f Makefile_docker golang-run-makefile-target

.PHONY: docker-update-random-package
docker-update-random-package: MAKEFILE_TARGET=update-random-package
docker-update-random-package: .docker-golang-run-makefile-target

.PHONY: docker-list-packages-with-newer-upstream-versions
docker-list-packages-with-newer-upstream-versions: MAKEFILE_TARGET=list-packages-witn-newer-upstream-versions
docker-update-random-package: .docker-golang-run-makefile-target

.PHONY: docker-all
docker-all: MAKEFILE_TARGET=all
docker-all: docker-golang-run-makefile-target
