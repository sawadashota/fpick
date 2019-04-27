.DEFAULT_GOAL := help

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

_VERSION = $(shell git describe --tags $(git rev-list --tags --max-count=1))
_BUILT_TIME = $(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
_COMMIT_HASH = $(shell git rev-parse HEAD)
build: ## build wedding command
	CGO_ENABLED=0 go build -ldflags "-s -X github.com/sawadashota/fpick/cmd.Version=${_VERSION} -X github.com/sawadashota/fpick/cmd.BuildTime=${_BUILT_TIME} -X github.com/sawadashota/fpick/cmd.GitHash=${_COMMIT_HASH}" -a -o ${GOPATH}/bin/fpick github.com/sawadashota/fpick/cmd/fpick

# https://postd.cc/auto-documented-makefile/
help: ## show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'