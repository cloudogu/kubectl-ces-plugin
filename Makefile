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

ENV_CGO_ENABLED=CGO_ENABLED=0
ENV_GOARCH=GOARCH=${GOARCH}

ENV_GOOS_LINUX=GOOS=linux
ENV_BINARY_LINUX=BINARY=${LINUX_TARGET}/${ARTIFACT_ID}
ENV_GOOS_WINDOWS=GOOS=windows
ENV_BINARY_WINDOWS=BINARY=${WINDOWS_TARGET}/${ARTIFACT_ID}.exe
ENV_GOOS_DARWIN=GOOS=darwin
ENV_BINARY_DARWIN=BINARY=${DARWIN_TARGET}/${ARTIFACT_ID}

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
	@make -e ${LINUX_TARGET}/${ARTIFACT_ID}
	@make -e ${WINDOWS_TARGET}/${ARTIFACT_ID}.exe
	@make -e ${DARWIN_TARGET}/${ARTIFACT_ID}

${LINUX_TARGET}/${ARTIFACT_ID}: $(SRC) $(LINUX_TARGET)
	GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_LINUX}" \
    		${ENV_CGO_ENABLED} \
    		${ENV_GOARCH} \
    		${ENV_GOOS_LINUX} \
    		${ENV_BINARY_LINUX} \
		make -e compile

${WINDOWS_TARGET}/${ARTIFACT_ID}.exe: $(SRC) $(WINDOWS_TARGET)
	GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_WINDOWS}" \
    		${ENV_CGO_ENABLED} \
    		${ENV_GOARCH} \
    		${ENV_GOOS_WINDOWS} \
    		${ENV_BINARY_WINDOWS} \
	  	make -e compile

${DARWIN_TARGET}/${ARTIFACT_ID}: $(SRC) $(DARWIN_TARGET)
	GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_DARWIN}" \
		${ENV_CGO_ENABLED} \
		${ENV_GOARCH} \
		${ENV_GOOS_DARWIN} \
		${ENV_BINARY_DARWIN} \
		make -e compile

.PHONY: create-krew-archives
create-krew-archives: ${LINUX_TARGET}/${ARTIFACT_ID} ${WINDOWS_TARGET}/${ARTIFACT_ID}.exe ${DARWIN_TARGET}/${ARTIFACT_ID}
	@cp LICENSE ${LINUX_TARGET}
	@cd ${LINUX_TARGET} && tar -czvf "${ARTIFACT_ID}_linux_${GOARCH}.tar.gz" *
	@cp LICENSE ${WINDOWS_TARGET}
	@cd ${WINDOWS_TARGET} && zip "${ARTIFACT_ID}_windows_${GOARCH}.zip" *
	@cp LICENSE ${DARWIN_TARGET}
	@cd ${DARWIN_TARGET} && tar -czvf "${ARTIFACT_ID}_darwin_${GOARCH}.tar.gz" *


.PHONY: update-krew-version
update-krew-version: ## Update the kubectl plugin manifest with the current artefact version.
	@yq -i ".spec.version |= \"${VERSION}\"" ${KREW_MANIFEST}

KREW_ARCHIVE_CHECKSUMS=target/linux/kubectl-ces_linux_amd64.tar.gz target/windows/kubectl-ces_windows_amd64.zip target/darwin/kubectl-ces_darwin_amd64.tar.gz

.PHONY: update-krew-checksums
update-krew-checksums:
	@echo "Generating Checksums"

	export SHA256_LINUX="$$(sha256sum target/linux/kubectl-ces_linux_amd64.tar.gz | awk {'print $$1'})" ; \
	yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"linux\").sha256=\"$${SHA256_LINUX}\"" deploy/krew/plugin.yaml

	export SHA256_WINDOWS="$$(sha256sum target/windows/kubectl-ces_windows_amd64.zip | awk {'print $$1'})" ; \
	yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"windows\").sha256=\"$${SHA256_WINDOWS}\"" deploy/krew/plugin.yaml

	export SHA256_DARWIN="$$(sha256sum target/darwin/kubectl-ces_darwin_amd64.tar.gz | awk {'print $$gen1'})" ; \
	yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"darwin\").sha256=\"$${SHA256_DARWIN}\"" deploy/krew/plugin.yaml


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