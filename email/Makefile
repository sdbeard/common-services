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

# ==================================================================================
# Build support

build: clean
	docker build \
	-f ./_zerf/docker/Dockerfile.amd64 \
	-t ${BINARY}:1.0.0 \
	--build-arg VCS_REF=`git rev-parse HEAD` \
	--build-arg BINARY=kp-email \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	--build-arg VERSION=1.0.0 \
	.
	docker image prune -f

dev: clean
	DOCKER_BUILDKIT=1 \
	docker build \
	-f ./_zerf/docker/Dockerfile.dev \
	--build-arg VCS_REF=`git rev-parse HEAD` \
	--build-arg BINARY=kp-email \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	--build-arg VERSION=1.0.0 \
	--output type=local,dest=build \
	.
	@mv ./build/kp-email .

clean:
	$(info Cleaning previous build...)
	@if [ -d ./build ]; then rm -rf ./build; fi
	@if [ -f ./kp-email ]; then rm -rf ./kp-email; fi

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
