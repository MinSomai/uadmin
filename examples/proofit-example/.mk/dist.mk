DOCKER_IMAGE?=proofit/proofit
DOCKER_TAG?=release
DOCKER_IMAGE_TAG?=test
DESTDIR?=$(shell pwd)

.PHONY: docker-image
docker-image: static
	cp ./main contrib/docker/proofit.$$(uname -m)
	cp ./configs/dev.yml contrib/docker/proofit-example.yml
	docker build -t ${DOCKER_IMAGE}:${DOCKER_IMAGE_TAG} --build-arg ARCH=$$(uname -m) -f contrib/docker/Dockerfile contrib/docker/