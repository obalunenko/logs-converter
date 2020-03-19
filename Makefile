NAME=logs-converter
BIN_DIR=./bin


TARGET_MAX_CHAR_NUM=20


## Show help
help:
	${call colored, help is running...}
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-$(TARGET_MAX_CHAR_NUM)s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)



## Compile app
compile:
	${call, compile is running...}
	./scripts/compile.sh
.PHONY: compile

## Cross compile
cross-compile:
	${call, compile is running...}
	./scripts/cross-compile.sh
.PHONY: cross-compile

## lint project
lint:
	${call, lint is running...}
	./scripts/run-linters.sh
.PHONY: lint

lint-ci:
	${call, lint_ci is running...}
	./scripts/run-linters-ci.sh
.PHONY: lint-ci

## format markdown files in project
pretty-markdown:
	find . -name '*.md' -not -wholename './vendor/*' | xargs prettier --write
.PHONY: pretty-markdown

## Test all packages
test:
	${call, test is running...}
	./scripts/run-tests.sh
.PHONY: test

## Test coverage
test-cover:
	${call, test-cover is running...}
	./scripts/coverage.sh
.PHONY: test-cover

new-version: lint test compile
	${call, new version is running...}
	./scripts/version.sh
.PHONY: new-version


## Release
release:
	${call, release is running...}
	./scripts/release.sh
.PHONY: release

## Fix imports sorting
imports:
	${call, sort and group imports...}
	./scripts/fix-imports.sh
.PHONY: imports

## dependencies - fetch all dependencies for sripts
dependencies:
	${call, dependensies is running...}
	./scripts/get-dependencies.sh
.PHONY: dependencies

## Docker compose up
docker-up:
	${call, docker is running...}
	docker-compose -f ./docker-compose.yml up -d --build

.PHONY: docker-up

## Docker compose down
docker-down:
	${call, docker is running...}
	docker-compose -f ./docker-compose.yml down --volumes

.PHONY: docker-down

## review code
review:
	${call, review is running...}
	reviewdog -reporter=github-pr-check

.PHONY: review

.DEFAULT_GOAL := test
