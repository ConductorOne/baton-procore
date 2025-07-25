GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BUILD_DIR = dist/${GOOS}_${GOARCH}
GENERATED_CONF := pkg/config/conf.gen.go

ifeq ($(GOOS),windows)
OUTPUT_PATH = ${BUILD_DIR}/baton-procore.exe
else
OUTPUT_PATH = ${BUILD_DIR}/baton-procore
endif

# Set the build tag conditionally based on ENABLE_LAMBDA
ifdef BATON_LAMBDA_SUPPORT
	BUILD_TAGS=-tags baton_lambda_support
else
	BUILD_TAGS=
endif

.PHONY: build
build: ${GENERATED_CONF}
	go build ${BUILD_TAGS} -o ${OUTPUT_PATH} ./cmd/baton-procore

$(GENERATED_CONF): pkg/config/config.go go.mod
	@echo "Generating $(GENERATED_CONF)..."
	go generate -tags=generate ./pkg/config

generate: $(GENERATED_CONF)

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor

.PHONY: add-dep
add-dep:
	go mod tidy -v
	go mod vendor

.PHONY: lint
lint: format
	golangci-lint run

.PHONY: format
format:
	@golangci-lint fmt
