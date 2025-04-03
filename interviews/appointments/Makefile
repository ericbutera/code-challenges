.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: run
run:
	@echo "Running..."
	@go run cmd/api/main.go

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
