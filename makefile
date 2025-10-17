LOG_DIRS = ./internal/services/kaas/tmp ./internal/services/dbaas/tmp
REPORT_DIR = $(shell pwd)/reports
OS := $(shell uname)

.PHONY: prepare-go prepare-env test-unit generate-reports global-coverage cobertura install clean

prepare-go:
	@go mod download

prepare-env:
	@./scripts/setup.sh "$(LOG_DIRS)" "$(REPORT_DIR)"

test-unit: prepare-go prepare-env
	@./scripts/run-tests.sh "$(REPORT_DIR)"

generate-reports: prepare-go prepare-env
	@./scripts/generate-reports.sh "$(REPORT_DIR)" "$(OS)"

global-coverage: prepare-go prepare-env
	@./scripts/global-coverage.sh "$(REPORT_DIR)" "$(OS)"

cobertura: prepare-go prepare-env
	@./scripts/generate-cobertura-reports.sh "$(REPORT_DIR)" "$(OS)"

install:
	@go build -cover -o ${HOME}/go/bin/terraform-provider-infomaniak main.go

clean:
	@rm -rf $(REPORT_DIR)
