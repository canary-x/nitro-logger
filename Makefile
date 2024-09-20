GO=go
GO_TARGETS=./cmd/... ./internal/...

.PHONY: help
help:
	@grep -E '^[a-zA-Z_\-\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build logger
	@${GO} build -o nitro-logger

.PHONY: start
start: stop ## start logger
	@./nitro-logger >/dev/null 2>&1 & disown

.PHONY: stop
stop: ## stop logger
	@ps aux | grep nitro-logger | grep -v grep | awk '{print $$2}' | xargs -r kill -9
