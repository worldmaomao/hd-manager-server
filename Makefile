.PHONY: build clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/hd-manager

.PHONY: $(MICROSERVICES)

DOCKERS=docker_hd_manager
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
BASE_VERSION=$(shell cat ./BASE_VERSION 2>/dev/null || echo 0.0.0)


GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)

cmd/hd-manager:
	$(GO) build -mod=vendor -o $@ ./cmd


clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker.base:
	docker build -f ./Dockerfile.base \
		-t worldmaomao/hd-manager-base:$(BASE_VERSION) \
		.

docker.base.push:
	docker push worldmaomao/hd-manager-base:$(BASE_VERSION)


docker_hd_manager:
	docker build -f ./Dockerfile \
		-t worldmaomao/hd-manager:$(VERSION) \
		.

docker.push:
	docker push worldmaomao/hd-manager:$(VERSION)