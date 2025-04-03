VERSION=latest
MAIN=main.go
IMAGE_NAME=project
# IMAGE_REPO=ghcr.io/ericbutera
IMAGE_TAG=${VERSION}

.DEFAULT_GOAL := help

.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 GOOS=linux go build -o build/app

# .PHONY: image-build
# image-build: ## Build docker image
# 	docker build -t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG} .
# .PHONY: image-push
#  image-push: ## Push docker image
# 	docker push ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: compose-up
compose-up: ## Run docker compose up
	docker compose up --build

# TODO: compose-down

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: run
run:
	@echo "Running..."
	@go run cmd/api/main.go

.PHONY: lint
lint:
	golangci-lint run --fix --config .golangci.yaml --timeout 5m --concurrency 4

.PHONY: setup-asdf # install binaries from .tool-versions
setup-asdf:
	cut -d\  -f1 .tool-versions|grep -E '^[^#]'|xargs -L1 asdf plugin add
	asdf install

.PHONY: setup-pre-commit # install pre-commit hooks
setup-pre-commit:
	pre-commit install --install-hook
	pre-commit install --hook-type commit-msg

.PHONY: setup-dev-env # setup dev environment
setup-dev-env: setup-asdf setup-pre-commit


.PHONY: generate-mocks # mockery test mocks defined in .mockery.yaml
generate-mocks:
	mockery