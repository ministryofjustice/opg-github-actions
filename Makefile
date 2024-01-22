SHELL := $(shell which bash)
APPNAME := OPGGitHubActions

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
HOST_ARCH := ${OS}_${ARCH}

BUILD_FOLDER = ./builds/
OS_AND_ARCHS_TO_BUILD := darwin_arm64 darwin_amd64

PWD := $(shell pwd)
USER_PROFILE := ~/.zprofile


.DEFAULT_GOAL: self
.PHONY: self all requirements darwin_arm64 darwin_amd64 
.ONESHELL: self all requirements darwin_arm64 darwin_amd64 
.EXPORT_ALL_VARIABLES:

self: $(HOST_ARCH)


darwin_x86_64: requirements
	@${MAKE} darwin_amd64

darwin_arm64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)$@/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)$@/main main.go
	@echo Build $@ complete.

darwin_amd64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)$@/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)$@/main main.go
	@echo Build $@ complete.


requirements:
ifeq (, $(shell which go))
	$(error go command not found)
endif
ifndef GOBIN
	$(warning GOBIN is not defined, configuring as ${HOME}/go/bin)
	$(shell mkdir -p ${HOME}/go/bin )
	$(shell echo "" >> ${USER_PROFILE};)
	$(shell echo "# ADDED BY ${PWD}/Makefile" >> ${USER_PROFILE};)
	$(shell echo export GOBIN="\$${HOME}/go/bin" >> ${USER_PROFILE};)
	$(shell echo export PATH="\$${PATH}:\$${GOBIN}" >> ${USER_PROFILE})
endif
	@echo All requirements checked
	@rm -Rf ${BUILD_FOLDER}
	@test -f ${USER_PROFILE} && source ${USER_PROFILE} || echo ${USER_PROFILE} not found	
	@cd $(PWD)/go && go test -json ./... > ./test-results.json
	@echo All tests completed