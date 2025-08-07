CGO_ENABLED = 0 # statically linked = 0
TARGETOS=linux
ifeq ($(OS),Windows_NT) 
    TARGETOS := Windows
else
    TARGETOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown' | tr '[:upper:]' '[:lower:]')
endif
TARGETARCH = amd64

.DEFAULT_GOAL = help
PID = /tmp/serving.pid
APP_UID     ?= $(shell id -u)
DC = docker compose
#-----------------------------------------------------------------------------------------------------------------------

help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

.PHONY: up
up: ## Up all
	@$(DC) up -d

.PHONY: down
down: ## Down all
	@$(DC) down

.PHONY: shell
shell: ## Enter application go-client container
	@$(DC) exec -it go-client bash

chat: ## Start chat
	@$(DC) exec -it go-client ./bin/go-client chat

# https://ollama.com/library
.PHONY: get-models ## Download Ollama models
get-models:
#	@$(DC) exec -it ollama ollama pull llama3.1
#	@$(DC) exec -it ollama ollama pull gemma3
	@$(DC) exec -it ollama ollama pull qwen3


.PHONY: go-build
go-build: ## Build dev application (go build)
#    @$(DC) exec application sh -c "go mod tidy"
#    @$(DC) exec sh -c "env CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-X main.env=dev' -o bin/app ./"
	@$(DC) exec go-client sh -c "env CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/go-client ./"

.PHONY: clean
clean: ## Clean bin/
	@$(DC) exec dev sh -c "rm -rf bin/go-client"
