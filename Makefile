PKG := .
CMD := $(PKG)/cmd/schema-generate
BIN := schema-generate

# Build
SHELLCHECK_VERSION = v0.9.0
SHELLCHECK_URL = https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).linux.$(shell uname -m).tar.xz
SHELLCHECK_BIN = tools/bin/shellcheck-$(SHELLCHECK_VERSION)

$(SHELLCHECK_BIN):
	mkdir -p tools/bin
	if (which curl > /dev/null) then \
            curl -sSfL $(SHELLCHECK_URL); \
	else \
            wget -O- -nv $(SHELLCHECK_URL); \
    fi \
    | tar xJ -C tools/bin --strip-components=1 shellcheck-$(SHELLCHECK_VERSION)/shellcheck \
    && mv tools/bin/shellcheck $@

.PHONY: all clean

all: clean $(BIN)

$(BIN): generator.go jsonschema.go cmd/schema-generate/main.go
	@echo "+ Building $@"
	CGO_ENABLED="0" go build -v -o $@ $(CMD)

clean:
	@echo "+ Cleaning $(PKG)"
	go clean -i $(PKG)/...
	rm -f $(BIN)
	rm -rf test/*_gen

# Test

# generate sources
JSON := $(wildcard test/*.json)
GENERATED_SOURCE := $(patsubst %.json,%_gen/generated.go,$(JSON))
test/%_gen/generated.go: test/%.json
	@echo "\n+ Generating code for $@"
	@D=$(shell echo $^ | sed 's/.json/_gen/'); \
	[ ! -d $$D ] && mkdir -p $$D || true
	./schema-generate -o $@ -p $(shell echo $^ | sed 's/test\///; s/.json//')  $^

.PHONY: test codecheck fmt lint vet

test: $(BIN) $(GENERATED_SOURCE)
	@echo "\n+ Executing tests for $(PKG)"
	go test -v -race -cover $(PKG)/...


codecheck: fmt lint vet

fmt:
	@echo "+ go fmt"
	go fmt $(PKG)/...

lint: $(GOPATH)/bin/golint $(SHELLCHECK_BIN)
	@echo "+ go lint"
	golint -min_confidence=0.1 $(PKG)/...
	$(SHELLCHECK_BIN) ./test/gen.sh

$(GOPATH)/bin/golint:
	go get -v golang.org/x/lint/golint

vet:
	@echo "+ go vet"
	go vet $(PKG)/...
