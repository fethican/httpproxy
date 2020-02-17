HASH:=$(shell git rev-parse --short HEAD)
DIRTY:=$(shell bash -c 'if [ -n "$$(git status --porcelain --untracked-files=no)" ]; then echo -dirty; fi')
COMMIT ?= $(HASH)$(DIRTY)
Built:=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS:=-X main.Commit=$(COMMIT) -X main.BuiltAt=$(Built)

binary:
	CGO_ENABLED=0 go build -o /binary -ldflags "$(LDFLAGS)"
