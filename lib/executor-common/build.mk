# executor-common library build configuration

EXECUTOR_COMMON_LIB_NAME := executor-common
EXECUTOR_COMMON_LIB_DIR := lib/$(EXECUTOR_COMMON_LIB_NAME)
EXECUTOR_COMMON_OUT := $(OUT)/$(EXECUTOR_COMMON_LIB_NAME)

# Library-specific variables
EXECUTOR_COMMON_VERSION := $(shell cat version.txt)
EXECUTOR_COMMON_WHEEL_NAME := executor_common-$(EXECUTOR_COMMON_VERSION)-py3-none-any.whl
ARK_EXECUTOR_COMMON_WHL := $(EXECUTOR_COMMON_OUT)/dist/$(EXECUTOR_COMMON_WHEEL_NAME)

# Pre-calculate all stamp paths
EXECUTOR_COMMON_STAMP_DEPS := $(EXECUTOR_COMMON_OUT)/stamp-deps
EXECUTOR_COMMON_STAMP_BUILD := $(EXECUTOR_COMMON_OUT)/stamp-build
EXECUTOR_COMMON_STAMP_TEST := $(EXECUTOR_COMMON_OUT)/stamp-test

# Add library output directory to clean targets
CLEAN_TARGETS += $(EXECUTOR_COMMON_OUT)
# Clean up Python artifacts
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/__pycache__
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/.pytest_cache
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/.ruff_cache
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/*.egg-info
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/dist
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/build
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/.coverage
CLEAN_TARGETS += $(EXECUTOR_COMMON_LIB_DIR)/htmlcov

# Define phony targets
.PHONY: $(EXECUTOR_COMMON_LIB_NAME)-build $(EXECUTOR_COMMON_LIB_NAME)-test

# Dependencies
$(EXECUTOR_COMMON_LIB_NAME)-deps: $(EXECUTOR_COMMON_STAMP_DEPS)
$(EXECUTOR_COMMON_STAMP_DEPS): $(EXECUTOR_COMMON_LIB_DIR)/pyproject.toml | $(OUT)
	@mkdir -p $(dir $@)
	cd $(EXECUTOR_COMMON_LIB_DIR) && uv sync --all-extras
	@touch $@

# Build target
$(EXECUTOR_COMMON_LIB_NAME)-build: $(EXECUTOR_COMMON_STAMP_BUILD) # HELP: Build executor-common Python library wheel
$(EXECUTOR_COMMON_STAMP_BUILD): $(ARK_EXECUTOR_COMMON_WHL)
	@touch $@

# Build Python wheel
$(ARK_EXECUTOR_COMMON_WHL): $(EXECUTOR_COMMON_STAMP_DEPS) $(EXECUTOR_COMMON_LIB_DIR)/src/executor_common/*.py | $(OUT)
	@mkdir -p $(EXECUTOR_COMMON_OUT)
	cd $(EXECUTOR_COMMON_LIB_DIR) && uv run python -m build --outdir $(EXECUTOR_COMMON_OUT)/dist

# Test target
$(EXECUTOR_COMMON_LIB_NAME)-test: $(EXECUTOR_COMMON_STAMP_TEST) # HELP: Run executor-common tests
$(EXECUTOR_COMMON_STAMP_TEST): $(EXECUTOR_COMMON_STAMP_DEPS)
	@if [ -d "$(EXECUTOR_COMMON_LIB_DIR)/tests" ]; then \
		cd $(EXECUTOR_COMMON_LIB_DIR) && uv run python -m pytest tests/; \
	else \
		echo "No tests directory found for executor-common"; \
	fi
	@touch $@

