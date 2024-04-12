BINARY_NAME ?= myapp #change me
VERSION ?= 0.0.1
BUILD_DIR := bin
GIT_REV_PARSE := $(shell git rev-parse HEAD)
COMMIT_ID := $(if ${GIT_REV_PARSE},${GIT_REV_PARSE},unknown)
DATECMD := date$(if $(findstring Windows,$(OS)),.exe,)
BUILD_TIMESTAMP := $(shell ${DATECMD} +%Y-%m-%dT%H:%m:%S%z)
.DEFAULT_GOAL := all


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


.PHONEY: hello
## hello:
hello: ## dummy help build for the makefile
	echo "hello"
	echo "always print hello does not exist"

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
clean: ## clean the project	add "-" before command so if error it is ignored
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
test: ## run tests
ifeq (${MAKECMDGOALS},ci)
	go test -v -coverpkg=./... -coverprofile=coverage.out -json ./... -count=1 > test-report.json || (cat test-report.json; exit 1)
else
	go test -v -coverpkg=./... -coverprofile=coverage.out ./... -count=1
endif

.PHONY: cover
## cover:
cover: test ## run test coverage
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