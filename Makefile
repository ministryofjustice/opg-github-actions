SHELL := $(shell which bash)
APPNAME := OPGGitHubActions

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
HOST_ARCH := ${OS}_${ARCH}

BUILD_FOLDER = ./builds
OS_AND_ARCHS_TO_BUILD := darwin_arm64 darwin_amd64

PWD := $(shell pwd)
USER_PROFILE := ~/.zprofile


.DEFAULT_GOAL: all
.PHONY: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests build
.ONESHELL: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests build
.EXPORT_ALL_VARIABLES:

# when running self, requires you run a target based on your arch
all: $(HOST_ARCH)

# for the github action builder - so dont run requirements
linux_x86_64:
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@tar -cvf $(BUILD_FOLDER)/$@.tar $(BUILD_FOLDER)/$@
	@echo Build $@ complete.

# LOCAL DEV VERSIONS
darwin_x86_64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@cd $(PWD)/go && tar -cvf $(BUILD_FOLDER)/$@.tar $(BUILD_FOLDER)/$@
	@echo Build $@ complete.

darwin_arm64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@cd $(PWD)/go && tar -cvf $(BUILD_FOLDER)/$@.tar $(BUILD_FOLDER)/$@
	@echo Build $@ complete.

darwin_amd64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@cd $(PWD)/go && tar -cvf $(BUILD_FOLDER)/$@.tar $(BUILD_FOLDER)/$@
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
