PROJECT_NAME := Pulumi pde Resource Provider

PACK             := pde
PACKDIR          := sdk
PROJECT          := github.com/corymhall/pulumi-provider-pde
NODE_MODULE_NAME := @corymhall/pde
NUGET_PKG_NAME   := Corymhall.pde

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
COMPONENT_PATH   := components
VERSION_PATH     := ${PROVIDER_PATH}/cmd/main.Version
SCHEMA_FILE     := provider/cmd/pulumi-resource-pde/schema.json
C_SCHEMA_FILE     := components/cmd/pulumi-resource-pdec/schema.json

GOPATH			:= $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

build:: schema provider component build_go

component::
	rm -rf ${WORKING_DIR}/bin/${PROVIDER}c
	(cd components/cmd/${PROVIDER}c && VERSION=$(VERSION) SCHEMA=$(WORKING_DIR)/$(C_SCHEMA_FILE) go run generate.go)
	(cd components && go build -o $(WORKING_DIR)/bin/${PROVIDER}c -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${COMPONENT_PATH}/cmd/$(PROVIDER)c)

provider::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

install::
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin
	cp $(WORKING_DIR)/bin/${PROVIDER}c ${GOPATH}/bin

schema::
	(cd provider/cmd/$(CODEGEN) && go run main.go schema ../$(PROVIDER))
	(cd components/cmd/$(CODEGEN)c && go run main.go schema ../$(PROVIDER)c)

# provider:: bin/${PROVIDER}

ensure::
	cd components && go mod tidy
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd provider/tests && go mod tidy
	cd examples/simple-go && go mod tidy

build_go:: VERSION := $(shell pulumictl get version --language generic)
build_go:: schema
				rm -rf sdk/go
				cd provider/cmd/$(CODEGEN) && go run main.go go ../../../sdk/go ../$(PROVIDER)/schema.json $(VERSION)
				cd components/cmd/$(CODEGEN)c && go run main.go go ../../../sdk/go ../$(PROVIDER)c/schema.json $(VERSION)

generate_schema:: schema
