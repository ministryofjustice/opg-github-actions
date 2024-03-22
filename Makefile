SHELL := $(shell which bash)
APPNAME := OPGGitHubActions

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
HOST_ARCH := ${OS}_${ARCH}

BUILD_FOLDER = ./builds

PWD := $(shell pwd)
USER_PROFILE := ~/.zprofile


.DEFAULT_GOAL: all
.PHONY: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests test release test_release_notes test_release_download_self_build test_release_download_binary
.ONESHELL: all requirements darwin_arm64 darwin_amd64 linux_x86_64 tests test release test_release_notes test_release_download_self_build test_release_download_binary
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


test:
	@cd $(PWD)/go && env LOG_LEVEL="warn" LOG_TO="stdout" go test -v -run="$(test-name)"

tests:
	@cd $(PWD)/go && env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./...


# this checks release note generation for multi line echos to github_output work correctly
test_release_notes:
	@cd $(PWD)/.github/actions/; env DEBUG="1" TAGNAME="v3.0.0" GH_REPO="ministryofjustice/opg-github-actions" GH_COMMIT="main" LAST_TAG="v2.7.3" ./release-notes.sh

# this checks the self build version returns correctly with file path
test_release_download_self_build:
	@mkdir -p $(PWD)/test-release-download-self/tmp/
	@cd $(PWD)/.github/actions/; env DEBUG="1" GH_WORKSPACE="$(PWD)/test-release-download-self/tmp" GH_ACTION_REPOSITORY="ministryofjustice/opg-github-actions" GH_ACTION_REF="refs/heads/main" SELF="true" ./release-download.sh

# this test overwrites checks we can find a prebuilt version of the linux release binary
test_release_download_binary:
	@mkdir -p $(PWD)/test-release-download-bin/tmp
	@mkdir -p $(PWD)/test-release-download-bin/opg-gha
	@mkdir -p $(PWD)/test-release-download-bin/opg-gha-build
	@cd $(PWD)/.github/actions/; env DEBUG="1" hostBuild="linux_x86_64" GH_WORKSPACE="$(PWD)/test-release-download-bin/tmp" GH_ACTION_REPOSITORY="ministryofjustice/opg-github-actions" GH_ACTION_REF="v3.0.1-releasenotem.1" SELF="true" ./release-download.sh
	@cd $(PWD)/.github/actions/; env DEBUG="1" hostBuild="linux_x86_64" GH_WORKSPACE="$(PWD)/test-release-download-bin/tmp" GH_ACTION_REPOSITORY="ministryofjustice/opg-github-actions" GH_ACTION_REF="v3.0.2" SELF="true" ./release-download.sh