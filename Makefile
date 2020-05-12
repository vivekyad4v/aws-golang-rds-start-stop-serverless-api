ORG_ID ?= my
ENVIRON ?= local
PROJECT_NAME ?= $(notdir $(CURDIR))
PROJECT_ID ?= $(ORG_ID)-$(ENVIRON)-$(PROJECT_NAME)
AWS_REGION ?= $(AWS_DEFAULT_REGION)
AWS_BUCKET_NAME ?= $(PROJECT_ID)-artifacts-$(AWS_REGION)

export GO111MODULE=on

###### Sample for environment segregation for local development
ifeq ($(ENVIRON), uat)
	BUILD=umake
endif
ifeq ($(ENVIRON), prd)
	BUILD=pmake
endif

include .make

build:
	@ cd ./src/rdst ; go mod download ; GOOS=linux go build -o ./bin/rdst ./api/
	@ cd ./src/rdst ; go mod download ; GOOS=linux go build -o ./bin/auth ./auth/
