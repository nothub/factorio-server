LDFLAGS  := -extldflags=-static
GOFLAGS  := -tags netgo,timetzdata
GOFLAGS  += -ldflags="$(LDFLAGS)"

$(shell basename $(shell go list -m)): go.mod go.sum $(shell find * -type f -iname '*.go')
	@echo building $@ from $?
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -race -o $@ ./cmd
