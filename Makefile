# note: call scripts from /scripts

.PHONY: default build builder-image binary-image test stop clean-images clean push apply deploy

BUILDER = chowkidar-builder
BINARY = stakater/chowkidar
# TODO: we should be able specify the tag in command; so, that we can reuse same Makefile in Jenkins pipeline as well.
TAG = dev
REPOSITORY = ${BINARY}:${TAG}

VERSION=
BUILD=

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build test

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

builder-image:
	@docker build -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --rm "${BUILDER}" | docker build -t "${REPOSITORY}" -f Dockerfile.run -

test:
	"$(GOCMD)" test -race -v $(shell go list ./... | grep -v '/vendor/')

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i

push: ## push the latest Docker image to DockerHub
	docker push $(REPOSITORY)

apply:
	kubectl apply -f deployments/manifests/

deploy: binary-image push apply