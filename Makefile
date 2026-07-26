# go-for-systems-engineers — root Makefile.
#
# One interface for five toolchains (Go, Rust, OCaml, Swift, Kotlin). Every
# target works from the repo root and can be scoped to a single module with
# `make <target> M=03`. CI invokes these same targets, so local and CI runs are
# identical.

SHELL := bash
.DEFAULT_GOAL := help

# --- Toolchain versions (override to test another; CI passes these in) ---------
GO_VERSION      ?= 1.26.5
RUST_VERSION    ?= 1.92.0
OCAML_SWITCH    ?= 5.4.0
SWIFT_VERSION   ?= 6.2.3
KOTLIN_VERSION  ?= 2.4.10

export GOTOOLCHAIN := auto
export PATH := $(shell go env GOPATH)/bin:$(HOME)/.cargo/bin:$(PATH)

# Homebrew's OpenJDK is keg-only; wire it in on macOS only. On CI (Linux),
# setup-java provides java/JAVA_HOME and this block is skipped.
ifneq ($(wildcard /opt/homebrew/opt/openjdk/bin),)
  export PATH := /opt/homebrew/opt/openjdk/bin:$(PATH)
  export JAVA_HOME := /opt/homebrew/opt/openjdk
endif

# Run a command inside the pinned opam switch (so dune/ocaml resolve). CI (which
# uses setup-ocaml's own switch) overrides this with OPAM="opam exec --".
OPAM ?= opam exec --switch=$(OCAML_SWITCH) --

CURDIR := $(shell pwd)
CAPTURE := $(CURDIR)/bin/capture

# --- Module discovery ---------------------------------------------------------
# M scopes every target to modules/M-*. Empty M means all modules + tools.
ifeq ($(strip $(M)),)
  SCOPE := .
  MODULE_GLOB := modules
else
  SCOPE := $(wildcard modules/$(M)-*)
  MODULE_GLOB := $(SCOPE)
endif

# Dirs named reject-* hold code that is EXPECTED TO FAIL to compile (orphan-rule
# and exhaustiveness demos); capture builds them on purpose, the build/test/lint
# targets skip them. Reader exercises are red by design; solutions run separately.
REJECT := -not -path '*/reject-*/*'
GO_DIRS := $(shell find $(SCOPE) -name go.mod \
	-not -path '*/exercises/*' -not -path '*/solutions/*' $(REJECT) \
	2>/dev/null | xargs -n1 dirname | sort -u)
RUST_DIRS := $(shell find $(MODULE_GLOB) -name Cargo.toml -not -path '*/target/*' \
	$(REJECT) 2>/dev/null | xargs -n1 dirname | sort -u)
OCAML_DIRS := $(shell find $(MODULE_GLOB) -name dune-project -not -path '*/_build/*' \
	$(REJECT) 2>/dev/null | xargs -n1 dirname | sort -u)
SWIFT_DIRS := $(shell find $(MODULE_GLOB) -name Package.swift -not -path '*/.build/*' \
	$(REJECT) 2>/dev/null | xargs -n1 dirname | sort -u)
KOTLIN_DIRS := $(shell find $(MODULE_GLOB) -type d -name kotlin -not -path '*/build/*' \
	$(REJECT) 2>/dev/null | sort -u)
SOLUTION_DIRS := $(shell find $(MODULE_GLOB) -path '*/solutions/*' -name go.mod \
	2>/dev/null | xargs -n1 dirname | sort -u)

# --- Help ---------------------------------------------------------------------
.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
		{printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

# --- Setup --------------------------------------------------------------------
.PHONY: setup
setup: ## Print toolchain status (does not install; see README for installs)
	@bash scripts/versions.sh

.PHONY: tools
tools: ## Build the capture binary (falsifiability engine)
	@mkdir -p bin
	@go build -C tools/capture -o $(CAPTURE) .

# --- Build --------------------------------------------------------------------
.PHONY: build
build: build-go build-rust build-ocaml build-swift build-kotlin ## Build all

.PHONY: build-go
build-go: ## go build every module
	@set -e; for d in $(GO_DIRS); do echo "== go build $$d"; \
		(cd $$d && go build ./...); done

.PHONY: build-rust
build-rust: ## cargo build every crate
	@set -e; for d in $(RUST_DIRS); do echo "== cargo build $$d"; \
		(cd $$d && cargo build --quiet); done

.PHONY: build-ocaml
build-ocaml: ## dune build every project
	@set -e; for d in $(OCAML_DIRS); do echo "== dune build $$d"; \
		(cd $$d && $(OPAM) dune build); done

.PHONY: build-swift
build-swift: ## swift build every package
	@set -e; for d in $(SWIFT_DIRS); do echo "== swift build $$d"; \
		(cd $$d && swift build); done

.PHONY: build-kotlin
build-kotlin: ## kotlinc-compile every kotlin module (lib+app -> demo.jar)
	@set -e; for d in $(KOTLIN_DIRS); do echo "== kotlinc $$d"; \
		(cd $$d && mkdir -p build && \
		 kotlinc lib app -include-runtime -d build/demo.jar); done

# --- Test ---------------------------------------------------------------------
.PHONY: test
test: test-go test-rust test-ocaml test-swift test-kotlin ## Test all

.PHONY: test-go
test-go: ## go test every module
	@set -e; for d in $(GO_DIRS); do echo "== go test $$d"; \
		(cd $$d && go test ./...); done

.PHONY: test-rust
test-rust: ## cargo test every crate
	@set -e; for d in $(RUST_DIRS); do echo "== cargo test $$d"; \
		(cd $$d && cargo test --quiet); done

.PHONY: test-ocaml
test-ocaml: ## dune runtest every project
	@set -e; for d in $(OCAML_DIRS); do echo "== dune runtest $$d"; \
		(cd $$d && $(OPAM) dune runtest); done

.PHONY: test-swift
test-swift: ## swift test every package
	@set -e; for d in $(SWIFT_DIRS); do echo "== swift test $$d"; \
		(cd $$d && swift test); done

.PHONY: test-kotlin
test-kotlin: ## kotlinc-compile lib+test and run the assertion main
	@set -e; for d in $(KOTLIN_DIRS); do echo "== kotlin test $$d"; \
		bash scripts/test-kotlin.sh $$d; done

.PHONY: test-race
test-race: ## go test -race -count=2 every module
	@set -e; for d in $(GO_DIRS); do echo "== go test -race $$d"; \
		(cd $$d && go test -race -count=2 ./...); done

.PHONY: solutions
solutions: ## Verify exercise corrigés pass (excluded from default test)
	@set -e; for d in $(SOLUTION_DIRS); do echo "== solutions $$d"; \
		(cd $$d && go test ./...); done

.PHONY: exercises
exercises: ## Run reader exercises (RED by design until you solve them)
	@for d in $$(find $(MODULE_GLOB) -path '*/exercises/*' -name go.mod \
		2>/dev/null | xargs -n1 dirname | sort -u); do \
		echo "== exercises $$d (expected to fail)"; \
		(cd $$d && go test ./...) || true; done

# --- Lint / format ------------------------------------------------------------
.PHONY: lint
lint: lint-go lint-rust lint-ocaml lint-swift lint-kotlin lint-shell ## Lint all

.PHONY: lint-go
lint-go: ## gofmt check, go vet, staticcheck, golangci-lint
	@set -e; test -z "$$(gofmt -l $$(find $(SCOPE) -name '*.go' \
		-not -path '*/_build/*' $(REJECT)))" || \
		{ echo "gofmt: files need formatting"; \
		  gofmt -l $$(find $(SCOPE) -name '*.go' $(REJECT)); exit 1; }
	@set -e; for d in $(GO_DIRS); do echo "== vet/staticcheck $$d"; \
		(cd $$d && go vet ./... && staticcheck ./... && \
		 golangci-lint run ./... --config $(CURDIR)/.golangci.yml); done

.PHONY: lint-rust
lint-rust: ## cargo fmt --check + clippy -D warnings
	@set -e; for d in $(RUST_DIRS); do echo "== fmt/clippy $$d"; \
		(cd $$d && cargo fmt --check && \
		 cargo clippy --all-targets --quiet -- -D warnings); done

.PHONY: lint-ocaml
lint-ocaml: ## ocamlformat check + dune build (warnings-as-errors)
	@set -e; for d in $(OCAML_DIRS); do echo "== ocamlformat/build $$d"; \
		(cd $$d && $(OPAM) dune build @fmt && $(OPAM) dune build); done

.PHONY: lint-swift
lint-swift: ## swift format lint (strict) + swift build (type check)
	@set -e; for d in $(SWIFT_DIRS); do echo "== swift-format/build $$d"; \
		(cd $$d && swift format lint --strict --recursive Sources Tests && \
		 swift build); done

.PHONY: lint-kotlin
lint-kotlin: ## ktlint + kotlinc -Werror
	@set -e; for d in $(KOTLIN_DIRS); do echo "== ktlint $$d"; \
		(cd $$d && ktlint "lib/**/*.kt" "app/**/*.kt" "test/**/*.kt" && \
		 mkdir -p build && \
		 kotlinc -Werror lib app -include-runtime -d build/demo.jar); done

.PHONY: lint-shell
lint-shell: ## shellcheck the harness + CI scripts
	@if command -v shellcheck >/dev/null; then \
		shellcheck scripts/*.sh scripts/ci/*.sh; \
	else echo "shellcheck not installed; skipping"; fi

.PHONY: format
format: ## Format all code in place
	@gofmt -w $$(find $(SCOPE) -name '*.go' -not -path '*/_build/*' $(REJECT))
	@for d in $(RUST_DIRS); do (cd $$d && cargo fmt); done
	@for d in $(OCAML_DIRS); do (cd $$d && $(OPAM) dune build @fmt \
		--auto-promote 2>/dev/null || true); done
	@for d in $(SWIFT_DIRS); do (cd $$d && swift format --in-place \
		--recursive Sources Tests); done
	@for d in $(KOTLIN_DIRS); do (cd $$d && ktlint -F \
		"lib/**/*.kt" "app/**/*.kt" "test/**/*.kt" || true); done

.PHONY: exhaustive
exhaustive: ## Run Module 02's sum-type analyzer on the exhaustive packages
	@d=modules/02-sum-types/go; \
	if [ -d $$d ]; then echo "== exhaustive $$d"; \
		(cd $$d && go run ./cmd/exhaustive .); fi

.PHONY: vuln
vuln: ## govulncheck every Go module (reachability-aware)
	@set -e; for d in $(GO_DIRS); do echo "== govulncheck $$d"; \
		(cd $$d && govulncheck ./...); done

# --- Docs (falsifiability) ----------------------------------------------------
.PHONY: docs
docs: tools ## Regenerate all captured README blocks
	@$(OPAM) $(CAPTURE) -root $(SCOPE)

.PHONY: docs-check
docs-check: tools ## Fail if any portable captured block is stale (docs-fresh gate)
	@$(OPAM) $(CAPTURE) -root $(SCOPE) -check

# --- Bench (manual; shared runners are too noisy for CI) ----------------------
.PHONY: bench
bench: ## How to run benchmarks (see each module README for exact invocations)
	@echo "Benchmarks are per-module and run manually. Example:"
	@echo "  bash scripts/go-bench.sh modules/00-bootstrap/go \\"
	@echo "       BenchmarkSum 10 modules/00-bootstrap/go/bench.txt"

# --- Aggregate + housekeeping -------------------------------------------------
.PHONY: ci
ci: build test test-race lint solutions docs-check ## Everything CI runs

.PHONY: module
module: ## Build+test+lint+docs-check one module: make module M=03
	@$(MAKE) build test lint docs-check M=$(M)

.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf bin
	@for d in $(RUST_DIRS); do (cd $$d && cargo clean --quiet 2>/dev/null); done
	@for d in $(OCAML_DIRS); do (cd $$d && rm -rf _build); done
	@for d in $(SWIFT_DIRS); do (cd $$d && rm -rf .build); done
	@for d in $(KOTLIN_DIRS); do (cd $$d && rm -rf build); done
	@find . -name '*.raw.txt' -delete
