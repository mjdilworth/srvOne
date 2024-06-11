PACKAGE_PATH ?= .
BINARY_NAME ?= srvOne#change me
VERSION ?= 0.0.1
BUILD_DIR := bin
GIT_REV_PARSE := $(shell git rev-parse HEAD)
COMMIT_ID := $(if ${GIT_REV_PARSE},${GIT_REV_PARSE},unknown)
DATECMD := date$(if $(findstring Windows,$(OS)),.exe,)
BUILD_TIMESTAMP := $(shell ${DATECMD} +%Y-%m-%dT%H:%m:%S%z)
.DEFAULT_GOAL := all
CONFIG_PATH ?= main


os_arch = $(word 4, $(shell go version))
os = $(word 1,$(subst /, ,$(os_arch)))
arch = $(word 2,$(subst /, ,$(os_arch)))


DOCKER_REGISTRY := #if set it should finished by /
EXPORT_RESULT := false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONEY: all
## all:
all: ## default is to build just for current platform
	@make me

.PHONEY: run
## run:
run: ## quick go run cmd/cli/main.go
	go run ${PACKAGE_PATH}/main.go

.PHONEY: hello
## hello:
hello: ## dummy hello build for the makefile to dump
	@echo "hello and dump out stuff"
	@echo "always print hello does not exist"
	@echo ${os_arch}
	@echo ${os}
	@echo ${arch}
	@echo ${VERSION}
	@echo ${BUILD_DIR}
	@echo ${GIT_REV_PARSE}
	$(eval BINARY := ${BINARY_NAME}$(if $(findstring windows,$(GOOS)),.exe,))
	

.PHONY: confirm
## confirm:
confirm: ## used to ask for user input if they want to continue
	@echo '\nAre you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: build-platform
## build-platform:
build-platform: ## build project to platform
	@echo Building ${GOOS}-${GOARCH}
	$(eval BINARY := ${BINARY_NAME}$(if $(findstring windows,$(GOOS)),.exe,))
	go build -v -o ${BUILD_DIR}/${GOOS}-${GOARCH}/$(BINARY) \
	-ldflags=all="-s -w -X ${CONFIG_PATH}.Version=${VERSION} -X ${CONFIG_PATH}.CommitId=${COMMIT_ID} -X ${CONFIG_PATH}.BuildTimestamp=${BUILD_TIMESTAMP}" ${PACKAGE_PATH}


.PHONEY: me
## me:
me: ## build for local architrcrture
	@make --no-print-directory build-platform GOOS=${os} GOARCH=${arch}  CGO_ENABLED=0
	

.PHONEY: build
## build:
build: ## build your project 
	@make --no-print-directory build-platform GOOS=windows GOARCH=amd64 CGO_ENABLED=0
	@make --no-print-directory build-platform GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	@make --no-print-directory build-platform GOOS=linux GOARCH=arm64 CGO_ENABLED=0
	@make --no-print-directory build-platform GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
	@make --no-print-directory build-platform GOOS=${os} GOARCH=${arch}  CGO_ENABLED=0
	
.PHONY: win
## win:
 win: ## build for windows
	@make --no-print-directory build-platform GOOS=windows GOARCH=amd64 KERBEROS_DEFAULT=true

.PHONY: build_pi
## build_pi:
build_pi: ## build your project for an old pi - arch arm
	GOARCH=arm	GOARM=5	GOOS=linux go build -o ${BINARY_NAME}-pi ./${BUILD_DIR}/.

.PHONY: clean
## clean:
clean: confirm ## clean the project	add "-" before command so if error it is ignored
	go	clean
	rm -rf ${BUILD_DIR}

.PHONY: fmt
## fmt:
fmt: ## format nicely
	go fmt ./...

.PHONY: tidy
## tidy:
tidy: ## mod tidy
	go mod tidy -v

.PHONY: test
## test:
test: ## run tests with test failure ignored (-)
	-go test -v -race -coverpkg=./... -coverprofile=coverage.out ./... -count=1

.PHONY: cover
## cover:
cover: test ## run test coverage - output to html page
	go tool cover -html=coverage.out

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)