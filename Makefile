# Set these to the desired values
ARTIFACT_ID=kubectl-ces
VERSION=v0.1.0

GOTAG=1.20.1
MAKEFILES_VERSION=7.5.0

.DEFAULT_GOAL:=help

GOARCH?=amd64
#GO_ENV_VARS=GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0
GO_BUILD_FLAGS?=-a -tags netgo -ldflags "-extldflags -static -s -w -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID)" -installsuffix cgo
CUSTOM_GO_MOUNT?=-v /tmp:/tmp

GITHUB_DOWNLOAD_URI=https://github.com/cloudogu/kubectl-ces-plugin/releases/download

LINUX_TARGET=$(TARGET_DIR)/linux
WINDOWS_TARGET=$(TARGET_DIR)/windows
DARWIN_TARGET=$(TARGET_DIR)/darwin
LINUX_TARGET_DIRSTAMP=${LINUX_TARGET}/.dirstamp
WINDOWS_TARGET_DIRSTAMP=${WINDOWS_TARGET}/.dirstamp
DARWIN_TARGET_DIRSTAMP=${DARWIN_TARGET}/.dirstamp

ENV_CGO_ENABLED=CGO_ENABLED=0
ENV_GOARCH=GOARCH=${GOARCH}
ENV_GOOS_LINUX=GOOS=linux
ENV_GOOS_WINDOWS=GOOS=windows
ENV_GOOS_DARWIN=GOOS=darwin
ENV_BINARY_LINUX=BINARY=${BINARY_LINUX}
ENV_BINARY_WINDOWS=BINARY=${BINARY_WINDOWS}
ENV_BINARY_DARWIN=BINARY=${BINARY_DARWIN}

BINARY_LINUX=${LINUX_TARGET}/${ARTIFACT_ID}
BINARY_WINDOWS=${WINDOWS_TARGET}/${ARTIFACT_ID}.exe
BINARY_DARWIN=${DARWIN_TARGET}/${ARTIFACT_ID}
KREW_ARCHIVE_LINUX=${ARTIFACT_ID}_linux_${GOARCH}.tar.gz
KREW_ARCHIVE_WINDOWS=${ARTIFACT_ID}_windows_${GOARCH}.zip
KREW_ARCHIVE_DARWIN=${ARTIFACT_ID}_darwin_${GOARCH}.tar.gz

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
include build/make/digital-signature.mk
include build/make/release.mk

${LINUX_TARGET_DIRSTAMP}:
	@mkdir -p $(LINUX_TARGET)
	@touch $@

$(WINDOWS_TARGET_DIRSTAMP):
	@mkdir -p $(WINDOWS_TARGET)
	@touch $@

$(DARWIN_TARGET_DIRSTAMP):
	@mkdir -p $(DARWIN_TARGET)
	@touch $@

.PHONY: compile-crossplatform
compile-crossplatform:  ## Compile the plugin for linux, windows and darwin.
	@make -e ${LINUX_TARGET}/${ARTIFACT_ID}
	@make -e ${WINDOWS_TARGET}/${ARTIFACT_ID}.exe
	@make -e ${DARWIN_TARGET}/${ARTIFACT_ID}

${BINARY_LINUX}: $(SRC) $(LINUX_TARGET_DIRSTAMP)
	@echo "Compiling $@"
	@GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_LINUX}" \
    		${ENV_CGO_ENABLED} \
    		${ENV_GOARCH} \
    		${ENV_GOOS_LINUX} \
    		${ENV_BINARY_LINUX} \
		make -e compile
	@touch $@

${BINARY_WINDOWS}: $(SRC) $(WINDOWS_TARGET_DIRSTAMP)
	@echo "Compiling $@"
	@GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_WINDOWS}" \
    		${ENV_CGO_ENABLED} \
    		${ENV_GOARCH} \
    		${ENV_GOOS_WINDOWS} \
    		${ENV_BINARY_WINDOWS} \
	  	make -e compile
	@touch $@

${BINARY_DARWIN}: $(SRC) $(DARWIN_TARGET_DIRSTAMP)
	@echo "Compiling $@"
	@GO_ENV_VARS="${ENV_CGO_ENABLED} ${ENV_GOARCH} ${ENV_GOOS_DARWIN}" \
		${ENV_CGO_ENABLED} \
		${ENV_GOARCH} \
		${ENV_GOOS_DARWIN} \
		${ENV_BINARY_DARWIN} \
		make -e compile
	@touch $@

## Managing KREW plugin generation

.PHONY: krew-create-archives
krew-create-archives: ${LINUX_TARGET}/${KREW_ARCHIVE_LINUX} ${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS} ${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN} ## Create KREW archives for all supportedd OSs.

${LINUX_TARGET}/${KREW_ARCHIVE_LINUX}: ${BINARY_LINUX}
	@cp LICENSE ${LINUX_TARGET}
	@if test -f $@ ; then \
  		echo "Found existing KREW archive $@. Deleting..." ; \
  		rm $@ ; \
  		echo "Continue archiving..." ; \
  	fi
	@cd ${LINUX_TARGET} && tar -czvf ${KREW_ARCHIVE_LINUX} * > /dev/null

${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS}: ${BINARY_WINDOWS}
	@cp LICENSE ${WINDOWS_TARGET}
	@if test -f $@ ; then \
  		echo "Found existing KREW archive $@. Deleting..." ; \
  		rm $@ ; \
  		echo "Continue archiving..." ; \
  	fi
	@cd ${WINDOWS_TARGET} && zip ${KREW_ARCHIVE_WINDOWS} * > /dev/null

${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN}: ${BINARY_DARWIN}
	@cp LICENSE ${DARWIN_TARGET}
	@if test -f $@ ; then \
  		echo "Found existing KREW archive $@. Deleting..." ; \
  		rm $@ ; \
  		echo "Continue archiving..." ; \
  	fi
	@cd ${DARWIN_TARGET} && tar -czvf ${KREW_ARCHIVE_DARWIN} * > /dev/null

.PHONY: krew-update-manifest-versions
krew-update-manifest-versions: ## Update the kubectl plugin manifest with the current artefact version.
	@yq -i ".spec.version |= \"${VERSION}\"" ${KREW_MANIFEST}
	@yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"linux\").uri=\"${GITHUB_DOWNLOAD_URI}/${VERSION}/${KREW_ARCHIVE_LINUX}\"" deploy/krew/plugin.yaml
	@yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"windows\").uri=\"${GITHUB_DOWNLOAD_URI}/${VERSION}/${KREW_ARCHIVE_WINDOWS}\"" deploy/krew/plugin.yaml
	@yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"darwin\").uri=\"${GITHUB_DOWNLOAD_URI}/${VERSION}/${KREW_ARCHIVE_DARWIN}\"" deploy/krew/plugin.yaml

.PHONY: krew-update-checksums
krew-update-checksums: ${LINUX_TARGET}/${KREW_ARCHIVE_LINUX} ${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS} ${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN}  ## Update SHA256 checksums in the KREW manifest.
	@echo "Generating Checksums"

	@export SHA256_LINUX="$$(sha256sum ${LINUX_TARGET}/${KREW_ARCHIVE_LINUX} | awk {'print $$1'})" ; \
		# update field sha256 with checksum whose sibling element match the value "linux" \
		yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"linux\").sha256=\"$${SHA256_LINUX}\"" deploy/krew/plugin.yaml

	@export SHA256_WINDOWS="$$(sha256sum ${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS} | awk {'print $$1'})" ; \
		yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"windows\").sha256=\"$${SHA256_WINDOWS}\"" deploy/krew/plugin.yaml

	@export SHA256_DARWIN="$$(sha256sum ${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN} | awk {'print $$1'})" ; \
		yq -i ".spec.platforms[] |= select(.selector.matchLabels.os == \"darwin\").sha256=\"$${SHA256_DARWIN}\"" deploy/krew/plugin.yaml

.PHONY krew-collect:
krew-collect: ${LINUX_TARGET}/${KREW_ARCHIVE_LINUX} ${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS} ${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN} ## Move all cross-compiles KREW archives in the target directory.
	@echo "Moving archives to ${TARGET_DIR}"
	@mv ${LINUX_TARGET}/${KREW_ARCHIVE_LINUX} ${WINDOWS_TARGET}/${KREW_ARCHIVE_WINDOWS} ${DARWIN_TARGET}/${KREW_ARCHIVE_DARWIN} ${TARGET_DIR}

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
	@docker run --rm \
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
