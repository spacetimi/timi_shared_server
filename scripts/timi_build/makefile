GOPATH:=$(shell go env GOPATH)

TIMI_BUILD_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: timi_build

timi_build: FORCE
	@cd ${TIMI_BUILD_PATH}; go install -ldflags "-X main._META_executable_build_epoch_string=`date +%s`"

clean:
	@cd ${TIMI_BUILD_PATH}; go clean -i

FORCE:
