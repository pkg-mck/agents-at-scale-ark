# postgres-memory service build configuration

POSTGRES_MEMORY_SERVICE_NAME := postgres-memory
POSTGRES_MEMORY_SERVICE_DIR := services/$(POSTGRES_MEMORY_SERVICE_NAME)
POSTGRES_MEMORY_OUT := $(OUT)/$(POSTGRES_MEMORY_SERVICE_NAME)

# Service-specific variables
POSTGRES_MEM_IMAGE := postgres-memory
POSTGRES_MEM_TAG ?= latest
POSTGRES_MEM_HELM_RELEASE := postgres-memory

# Pre-calculate all stamp paths
POSTGRES_MEMORY_STAMP_TEST := $(POSTGRES_MEMORY_OUT)/stamp-test
POSTGRES_MEMORY_STAMP_BUILD := $(POSTGRES_MEMORY_OUT)/stamp-build
POSTGRES_MEMORY_STAMP_INSTALL := $(POSTGRES_MEMORY_OUT)/stamp-install

# Add service output directory to clean targets
CLEAN_TARGETS += $(POSTGRES_MEMORY_OUT)
# Clean up Go artifacts
CLEAN_TARGETS += $(POSTGRES_MEMORY_SERVICE_DIR)/cover.out
CLEAN_TARGETS += $(POSTGRES_MEMORY_SERVICE_DIR)/.cover

# Define phony targets
.PHONY: $(POSTGRES_MEMORY_SERVICE_NAME)-build $(POSTGRES_MEMORY_SERVICE_NAME)-install $(POSTGRES_MEMORY_SERVICE_NAME)-uninstall $(POSTGRES_MEMORY_SERVICE_NAME)-dev $(POSTGRES_MEMORY_SERVICE_NAME)-test

# Test target (Go service)
$(POSTGRES_MEMORY_SERVICE_NAME)-test: $(POSTGRES_MEMORY_STAMP_TEST)
$(POSTGRES_MEMORY_STAMP_TEST): $(POSTGRES_MEMORY_SERVICE_DIR)/go.mod $(POSTGRES_MEMORY_SERVICE_DIR)/go.sum | $(OUT)
	@mkdir -p $(dir $@)
	cd $(POSTGRES_MEMORY_SERVICE_DIR) && go test ./... -coverprofile cover.out
	@touch $@

# Build target
$(POSTGRES_MEMORY_SERVICE_NAME)-build: $(POSTGRES_MEMORY_STAMP_BUILD) # HELP: Build PostgreSQL memory service Docker image
$(POSTGRES_MEMORY_STAMP_BUILD): $(POSTGRES_MEMORY_STAMP_TEST)
	cd $(POSTGRES_MEMORY_SERVICE_DIR) && docker build -t $(POSTGRES_MEM_IMAGE):$(POSTGRES_MEM_TAG) .
	@touch $@

# Install target (includes PGO installation)
$(POSTGRES_MEMORY_SERVICE_NAME)-install: $(POSTGRES_MEMORY_STAMP_INSTALL) # HELP: Deploy PostgreSQL memory service to cluster
$(POSTGRES_MEMORY_STAMP_INSTALL): $(POSTGRES_MEMORY_STAMP_BUILD)
	cd $(POSTGRES_MEMORY_SERVICE_DIR) && ./build.sh $(POSTGRES_MEM_TAG) auto
	helm install pgo oci://registry.developers.crunchydata.com/crunchydata/pgo || true
	sleep 5
	kubectl wait --for=condition=available deployment/pgo --timeout=300s
	cd $(POSTGRES_MEMORY_SERVICE_DIR) && helm install $(POSTGRES_MEM_HELM_RELEASE) ./chart
	@touch $@

# Uninstall target
$(POSTGRES_MEMORY_SERVICE_NAME)-uninstall: # HELP: Remove PostgreSQL memory service from cluster
	helm uninstall $(POSTGRES_MEM_HELM_RELEASE) --ignore-not-found
	rm -f $(POSTGRES_MEMORY_STAMP_INSTALL)

# Dev target
$(POSTGRES_MEMORY_SERVICE_NAME)-dev:
	cd $(POSTGRES_MEMORY_SERVICE_DIR) && go run main.go
