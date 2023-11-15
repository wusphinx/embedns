.PHONY: lint
lint:
	golangci-lint run -v 

build:
	go build 

COMMIT_ID := $(shell git rev-parse --short HEAD)

img_repo ?= local

.PHONY: docker
docker: 
	docker build -t $(img_repo):$(COMMIT_ID) .
