# Set these to the desired values
ARTIFACT_ID=kubectl-ces
VERSION=v0.1.0

GOTAG=1.18.6
MAKEFILES_VERSION=7.5.0

.DEFAULT_GOAL:=help

GOARCH?=amd64
#GO_ENV_VARS=GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0
GO_BUILD_FLAGS?=-a -tags netgo -ldflags "-extldflags -static -s -w -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID)" -installsuffix cgo
CUSTOM_GO_MOUNT?=-v /tmp:/tmp

LINUX_TARGET=$(TARGET_DIR)/linux
WINDOWS_TARGET=$(TARGET_DIR)/windows
DARWIN_TARGET=$(TARGET_DIR)/darwin

KREW_MANIFEST=deploy/krew/plugin.yaml

include build/make/variables.mk
GOMODULES=on

# You may want to overwrite existing variables for target actions to fit into your project.

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/mocks.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/package-debian.mk
include build/make/deploy-debian.mk
include build/make/digital-signature.mk
include build/make/release.mk

$(LINUX_TARGET):
	@mkdir -p $(LINUX_TARGET)
$(WINDOWS_TARGET):
	@mkdir -p $(WINDOWS_TARGET)
$(DARWIN_TARGET):
	@mkdir -p $(DARWIN_TARGET)

.PHONY: compile-crossplatform
compile-crossplatform:  ## Compile the plugin for linux, windows and darwin.
	@make ${LINUX_TARGET}/${ARTIFACT_ID}
	@make ${WINDOWS_TARGET}/${ARTIFACT_ID}
	@make ${DARWIN_TARGET}/${ARTIFACT_ID}

$(LINUX_TARGET)/$(ARTIFACT_ID): $(SRC) $(LINUX_TARGET)
	GO_ENV_VARS="CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH}" \
		CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=${GOARCH} \
		BINARY=${LINUX_TARGET}/${ARTIFACT_ID} \
		make -e compile

$(WINDOWS_TARGET)/$(ARTIFACT_ID): $(SRC) $(WINDOWS_TARGET)
	GO_ENV_VARS="CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH}" \
	  	CGO_ENABLED=0 \
	  	GOOS=windows \
	  	GOARCH=${GOARCH} \
	  	BINARY=${WINDOWS_TARGET}/${ARTIFACT_ID}.exe \
	  	make -e compile

$(DARWIN_TARGET)/$(ARTIFACT_ID): $(SRC) $(DARWIN_TARGET)
	GO_ENV_VARS="CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH}" \
		CGO_ENABLED=0 \
		GOOS=darwin \
		GOARCH=${GOARCH} \
		BINARY=${DARWIN_TARGET}/${ARTIFACT_ID} \
		make -e compile

.PHONY: update-krew-version
update-krew-version: ## Update the kubectl plugin manifest with the current artefact version.
	@yq ".spec.version |= \"${VERSION}\"" ${KREW_MANIFEST} > ${KREW_MANIFEST}.tmp
	@cp ${KREW_MANIFEST}.tmp ${KREW_MANIFEST}
	@rm ${KREW_MANIFEST}.tmp

KREW_ARCHIVE_CHECKSUMS=target/kubectl-ces_linux_amd64.tar.gz target/kubectl-ces_windows_amd64.zip target/kubectl-ces_darwin_amd64.tar.gz

$(BIN_CHECKSUMS): $(TARGET_DIR)
	@echo "Generating Checksums"
	@cd $(TARGET_DIR); find . -maxdepth 1 -not -type d | egrep -v ".(sha256sum|asc)$$" | xargs shasum -a 256 > $$(basename $@)

##@ Compiling go software

.PHONY: compile
compile: $(BINARY) ## Compile the application

ifeq ($(ENVIRONMENT), ci)

$(BINARY): $(SRC) vendor $(PRE_COMPILE)
	@echo "Built on CI server"
	$(GO_ENV_VARS) go build $(GO_BUILD_FLAGS) -o $(BINARY)

else

$(BINARY): $(SRC) vendor $(PASSWD) $(ETCGROUP) $(HOME_DIR) $(PRE_COMPILE)
	@echo "Building locally (in Docker)"
	docker run --rm \
		-e GOOS="${GOOS}" \
		-e GOARCH="${GOARCH}" \
		-e BINARY="${BINARY}" \
		-e CGO_ENABLED=0 \
		-e GO_ENV_VARS="${GO_ENV_VARS}" \
		-u "$(UID_NR):$(GID_NR)" \
		-v $(PASSWD):/etc/passwd:ro \
		-v $(ETCGROUP):/etc/group:ro \
		-v $(HOME_DIR):/home/$(USER) \
		-v $(WORKDIR):/go/src/github.com/cloudogu/$(ARTIFACT_ID) \
		$(CUSTOM_GO_MOUNT) \
		-w /go/src/github.com/cloudogu/$(ARTIFACT_ID) \
		$(GOIMAGE):$(GOTAG) \
  go build $(GO_BUILD_FLAGS) -o $(BINARY)

endif