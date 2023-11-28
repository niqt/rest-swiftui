
# Go Checkup
GOPATH ?= $(shell go env GOPATH)
GO111MODULE:=auto
export GO111MODULE
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif
PATH := ${GOPATH}/bin:$(PATH)

GO = go

Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m➡\033[0m")

# Commands
.PHONY: all
all: | deps

.PHONY: deps
deps: ; $(info $(M) Installing project dependencies ...) @ ## Install project dependencies
	$Q $(GO) mod tidy
#	$Q $(GO) mod vendor


.PHONY: build
build: ; $(info $(M) Building executable...) @ ## Build program binary
	$Q mkdir -p bin
	$Q ret=0 && for d in $$($(GO) list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		$(GO) build  ${GCFLAGS} -o bin/$${b} $$d || ret=$$? ; \
		echo "$(M) Build: bin/$${b}" ; \
		echo "$(M) Done!" ; \
	done ; exit $$ret

.PHONY: run
run: ; $(info $(M) Running dev build (on the fly) ...) @ ## Run intermediate builds
	$Q $(GO) run -race .


generate: export GOFLAGS := -mod=vendor
generate:
	@go mod vendor
	@go generate ./... | true
	@rm -rf ./vendor
