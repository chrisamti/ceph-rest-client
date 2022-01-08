.PHONY: all release clean build json lint help

UNAME_M := $(shell uname -m)
RACE=
ifeq ($(UNAME_M),x86_64)
	RACE=-race
endif

DOCKER_REGISTRY=dock.unitycms.io/feed-short-term/sam

all: generate lint test build ## Test, lint check and build application

update-feed-packages: ## update feed packages
	go get -u github.com/DND-IT/feed-cbdockertester@master
	go get -u github.com/DND-IT/feed-storage/v2/storage@master
	go get -u github.com/DND-IT/looney@master

	go mod tidy
	go mod vendor

release: clean ## Build release version of application
	mkdir -p ./dist
	CGO_ENABLED=0 go build -mod vendor -o ./dist/sam


build: ## Build application
	CGO_ENABLED=0 go build -mod vendor .

lint: ## Lint the project
	golangci-lint --timeout 300s run ./...

generate: ## Generate mocks
	go generate -mod vendor ./...

test:
	go test -mod vendor $(RACE) -v ./...

docker-alpine: ## create docker images
	#  create arm64
	docker buildx build -f docker/alpine/Dockerfile --platform linux/arm64 --tag $(DOCKER_REGISTRY):alpine3.12-arm64 .
	# create amd64
	docker buildx build -f docker/alpine/Dockerfile --platform linux/amd64 --tag $(DOCKER_REGISTRY):alpine3.12-amd64 .
	# push arm64
	docker push $(DOCKER_REGISTRY):alpine3.12-arm64
	# push amd64
	docker push $(DOCKER_REGISTRY):alpine3.12-amd64
	# create multi arch manifest
	# docker manifest create $(DOCKER_REGISTRY):alpine3.12-manual --amend $(DOCKER_REGISTRY):alpine3.12-arm64 --amend $(DOCKER_REGISTRY):alpine3.12-amd64
	docker manifest create $(DOCKER_REGISTRY):alpine3.12-manual --amend $(DOCKER_REGISTRY):alpine3.12-amd64
	# push
	docker manifest push $(DOCKER_REGISTRY):alpine3.12-manual

help: ## Print all possible targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)