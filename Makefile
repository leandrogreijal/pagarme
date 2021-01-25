GOCMD=go
DEPCMD=dep
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=pagarme
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DIRECTORY=bin

all:  make build test build build-linux
dep:
	$(DEPCMD) ensure
build: dep test
	$(GOBUILD) -o ${BINARY_DIRECTORY}/$(BINARY_NAME) -v
build-linux: dep test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ${BINARY_DIRECTORY}/$(BINARY_UNIX) -v
test: clean
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f ${BINARY_DIRECTORY}/$(BINARY_NAME)
	rm -f ${BINARY_DIRECTORY}/$(BINARY_UNIX)
run: build
	./${BINARY_DIRECTORY}/$(BINARY_NAME)