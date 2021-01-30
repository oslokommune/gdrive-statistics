SHELL = bash
IMAGE = ghcr.io/oslokommune/gdrive-statistics

GO := $(shell command -v go 2> /dev/null)
ifndef GO
$(error go is required, please install)
endif

GOPATH			:= $(shell go env GOPATH)
GOBIN			?= $(GOPATH)/bin
GOFUMPT			:= $(GOBIN)/gofumpt
GOLANGCILINT   	:= $(GOBIN)/golangci-lint

PKGS  = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
FILES = $(shell find . -name '.?*' -prune -o -name vendor -prune -o -name '*.go' -print)

help: ## Print this menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: ## - Set up stuff you need to run locally
	@echo "Install mkcert and run:"
	@echo "mkcert localhost"

gofumpt: ## -
	$(GO) get -u mvdan.cc/gofumpt

fmt: gofumpt  ## -
	$(GO) fmt $(PKGS)
	$(GOFUMPT) -s -w $(FILES)

golangcilint:
	# To bump, simply change the version at the end to the desired version. The git sha here points to the newest commit
	# of the install script verified by our team located here: https://github.com/golangci/golangci-lint/blob/master/install.sh
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/b90551cdf9c6214075f2a40d1b5595c6b41ffff0/install.sh | sh -s -- -b ${GOBIN} v1.32.2

lint: golangcilint ## -
	$(GOLANGCILINT) run

build-docker: ## -
	docker build . -t $(IMAGE)

run: ## -
	go run *.go

push: build-docker ## -
	docker push $(IMAGE)

up: ## docker-compose up
	(docker-compose down || true) && \
	docker-compose up -d && \
	docker logs -f gdrive-statistics_app_1

