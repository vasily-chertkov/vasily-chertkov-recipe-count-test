.DEFAULT_GOAL := build

PROJECT_NAME			?= recipe-count-task
PROJECT_ROOT			:= github.com/hellofreshdevtests/$(PROJECT_NAME)
REGISTRY				?= vasily.chertkov

GO_VER=1.18.3
ALPINE_VER=3.15

IMAGE_VERSION			?= 1.0.0
BECOME					?=

IMAGE_NAME				:= $(PROJECT_NAME)
DOCKERFILE_PATH			:= $(PWD)/docker/Dockerfile
DOCKER_IMAGE			:= $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_VERSION)

SRCROOT_IN_CONTAINER	:= /go/src/$(PROJECT_ROOT)
DOCKER_RUNNER 			:= $(BECOME) docker run --rm -u `id -u`:`id -g` -it
DOCKER_RUNNER			+= -w $(SRCROOT_IN_CONTAINER) -v $(PWD):$(SRCROOT_IN_CONTAINER)


BASE_IMAGE				:= golang:$(GO_VER)-alpine$(ALPINE_VER)
DOCKER_BUILDER			:= $(DOCKER_RUNNER) $(BASE_IMAGE)


.PHONY: build
build:
	$(BECOME) docker build \
		--build-arg GO_VER=$(GO_VER) \
		--build-arg ALPINE_VER=$(ALPINE_VER) \
		--build-arg WORKDIR=$(SRCROOT_IN_CONTAINER) \
		-t $(DOCKER_IMAGE) -f $(DOCKERFILE_PATH) .

.PHONY: clean
clean:
	$(BECOME) docker rmi -f $(shell docker images -q $(DOCKER_IMAGE)) || true
	$(BECOME) docker image prune -f --filter label=stage=server-intermediate

.PHONY: run
run: build .guard-INPUT_PATH
	$(BECOME) docker run \
		--rm \
		-u `id -u`:`id -g` \
		-v $$(shell realpath $(INPUT_PATH)):/input.json \
		--name $(PROJECT_NAME) \
		-p 80:8080 \
		$(DOCKER_IMAGE)

.PHONY: fmt
fmt: 
	$(DOCKER_BUILDER) go fmt $$(go list ./...)

.PHONY: test
test:
	$(DOCKER_RUNNER) \
		-e CGO_ENABLED=0 \
		-e GOCACHE=/tmp/.cache \
		$(BASE_IMAGE) \
		go test -v -mod=vendor ./...

