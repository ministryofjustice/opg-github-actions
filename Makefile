SHELL := $(shell which bash)
APPNAME := OPGGitHubActions

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
HOST_ARCH := ${OS}_${ARCH}

BUILD_FOLDER = ./builds

PWD := $(shell pwd)
USER_PROFILE := ~/.zprofile


.DEFAULT_GOAL: all
.PHONY: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests release
.ONESHELL: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests release
.EXPORT_ALL_VARIABLES:

# when running , requires you run a target based on your arch
all: $(HOST_ARCH)

release: $(HOST_ARCH)
	@cd $(PWD)/go/$(BUILD_FOLDER) && tar -czvf release.tar.gz *
# for the github action builder - so dont run requirements
linux_x86_64:
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@echo Build $@ complete.

# LOCAL DEV VERSIONS
darwin_x86_64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@echo Build $@ complete.

darwin_arm64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
	@echo Build $@ complete.

darwin_amd64: requirements
	@cd $(PWD)/go && mkdir -p $(BUILD_FOLDER)/
	@cd $(PWD)/go && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o $(BUILD_FOLDER)/$@ main.go
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


tests:
	@cd $(PWD)/go && env LOG_LEVEL="error"  LOG_TO="stdout" go test -v ./...