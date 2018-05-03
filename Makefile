.PHONY: all
all: fmt \
	list-packages-with-newer-upstream-versions \
	update-random-package \
	autodeb-server \
	vet \
	lint \
	data

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile

GO_IMPORT_PATH := salsa.debian.org/autodeb-team/autodeb

list-packages-with-newer-upstream-versions: $(SOURCES) install
	go build -v -o list-packages-with-newer-upstream-versions ${GO_IMPORT_PATH}/cmd/list-packages-with-newer-upstream-versions

update-random-package: $(SOURCES) install
	go build -v -o update-random-package ${GO_IMPORT_PATH}/cmd/update-random-package


autodeb-server: $(SOURCES) install
	go build -v -o autodeb-server ${GO_IMPORT_PATH}/cmd/autodeb-server

data:
	mkdir data

.PHONY: test
test:
	go test -v ${GO_IMPORT_PATH}/...

.PHONY: vet
vet: install
	go vet ${GO_IMPORT_PATH}/...

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
	# Binaries
	rm -f list-packages-with-newer-upstream-versions
	rm -f update-random-package
	rm -rf update-random-package-output-*

	# Other things created by this Makefile
	rm -rf data
	rm -f dependency-graph.svg

	# stuff produced at runtime
	rm -f database.sqlite

##
## Misc
##

dependency-graph.svg: $(SOURCES)
	go get github.com/kisielk/godepgraph
	godepgraph -o $(GO_IMPORT_PATH) $(GO_IMPORT_PATH)/cmd/autodeb-server | dot -Tsvg > dependency-graph.svg
	# This would also work:
	#    go get github.com/davecheney/graphpkg
	#    graphpkg -stdout -match '$(GO_IMPORT_PATH)' $(GO_IMPORT_PATH)/cmd/autodeb-server > dependency-graph.svg
