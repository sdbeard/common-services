#***********************************************************************************
# The MIT License (MIT)
#
# Copyright (c) 2023 Sean Beard

# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in the
# Software without restriction, including without limitation the rights to use, copy,
# modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
# and to permit persons to whom the Software is furnished to do so, subject to the
# following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
# PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
# CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
# SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
# ***********************************************************************************
SHELL := /bin/bash
BINARY := files-service
VERSION := 0.0.1
REGISTRY := registry.testlab.local:5000/
DOCKER_REMOTE_HOST := "ssh://testlab"

# ==================================================================================
# Build support

all: image deploy

build: cleanbin cleandocker
	docker build \
	-f ./_res/docker/Dockerfile.amd64 \
	-t ${BINARY}:${VERSION} \
	--build-arg VCS_REF=`git rev-parse HEAD` \
	--build-arg BINARY=${BINARY} \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	--build-arg VERSION=${VERSION} \
	.
	docker image prune -f

dev: clean
	DOCKER_BUILDKIT=1 \
	docker build \
	-f ./_res/docker/Dockerfile.bin \
	--build-arg VCS_REF=`git rev-parse HEAD` \
	--build-arg BINARY=${BINARY} \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	--build-arg VERSION=${VERSION} \
	--output type=local,dest=. \
	.

image: build
	docker tag ${BINARY}:${VERSION} ${REGISTRY}${BINARY}:${VERSION}
	docker tag ${BINARY}:${VERSION} ${REGISTRY}${BINARY}
	docker push ${REGISTRY}${BINARY}:${VERSION}
	docker push ${REGISTRY}${BINARY}

deploy: cleandeploy
	DOCKER_HOST=${DOCKER_REMOTE_HOST} docker-compose -f ./_res/compose/docker-compose.testlab.yml up -d
	
clean: cleanbin cleandocker cleandeploy
	$(info Cleaned all processes...")

cleanbin:
	$(info Cleaning previous build...)
	@if [ -f ./${BINARY} ]; then rm -rf ./${BINARY}; fi

cleandocker: 
	$(info Cleaning previous docker build...)
	docker images --format '{{.Repository}}:{{.Tag}}' | grep ${BINARY} | xargs -I {} docker rmi {}

cleandeploy:
	DOCKER_HOST=${DOCKER_REMOTE_HOST} docker-compose -f ./_res/compose/docker-compose.testlab.yml down
	DOCKER_HOST=${DOCKER_REMOTE_HOST} docker images --format '{{.Repository}}:{{.Tag}}' | grep ${BINARY} | DOCKER_HOST=${DOCKER_REMOTE_HOST} xargs -I {} docker rmi {}

# ==================================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==================================================================================
