BUILDPATH := $(CURDIR)
BINDIR := $(BUILDPATH)/bin
GO := $(shell which go)
# For future
# GOPATH ?= $(shell go env GOPATH)
GOBUILD := $(GO) build
TARGET_OS :=\
	darwin\
	linux\

PLATFORMS :=\
	amd64\

BINARY := reglogin

ifeq ($(GOHOSTOS),)
GOHOSTOS := $(shell uname | tr A-Z a-z)
endif

ifeq ($(GOHOSTARCH),)
GOHOSTARCH := $(shell uname -m | sed 's/x86_64/amd64/; s/^..86$$/386/;')
endif

HOSTBIN := $(BINARY)-$(GOHOSTOS)-$(GOHOSTARCH)

all: build

format:
	gofmt -s -w .

build: builddir
	@echo "Building reglogin ..."
	@$(GOBUILD) -o $(BINDIR)/$(HOSTBIN) -v
	@cp $(BINDIR)/$(HOSTBIN) $(BINDIR)/$(BINARY)
	@echo "Build completed!"

builddir:
	@echo "Creating destination folder ..."
	@if [[ ! -d $(BINDIR) ]]; then mkdir -p $(BINDIR); fi

build_all: builddir
	@echo "Building for $(TARGET_OS)"
	$(foreach GOOS,$(TARGET_OS),\
	$(foreach ARCH,$(PLATFORMS),\
	$(shell env GOOS=$(GOOS) GOARCH=$(ARCH) \
	$(GOBUILD) -o $(BINDIR)/$(BINARY)-$(GOOS)-$(ARCH))))
	@echo "Build completed!"

clean:
	@echo "Cleanup started ..."
	@rm -rf $(BINDIR)/
	@echo "Cleanup completed!"

.PHONY: clean build_all build builddir all
