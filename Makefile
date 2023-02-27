# Set these to the desired values
ARTIFACT_ID=kubectl-ces-plugin
VERSION=0.0.1

GOTAG=1.18.6
MAKEFILES_VERSION=7.5.0

.DEFAULT_GOAL:=help

BINARY=

include build/make/variables.mk
GOMODULES=on

# You may want to overwrite existing variables for target actions to fit into your project.

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
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