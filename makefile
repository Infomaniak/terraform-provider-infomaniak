LOG_DIRS = ./internal/services/kaas/tmp ./internal/services/dbaas/tmp
REPORT_DIR = $(shell pwd)/reports
OS := $(shell uname)


prepare-go:
	@go mod download

prepare-logs:
	@mkdir -p $(LOG_DIRS) 2>/dev/null || true
	@for dir in $(LOG_DIRS); do \
		touch $$dir/terraform.log; \
	done

prepare-reports:
	@mkdir -p $(REPORT_DIR)
	@mkdir -p $(REPORT_DIR)/unit
	@mkdir -p $(REPORT_DIR)/acceptance
	@mkdir -p $(REPORT_DIR)/covcounters/merged

testacc-unit: prepare-go prepare-logs prepare-reports
	@echo "=== Running acceptance tests ==="
	MOCKED=true GOCOVERDIR=$(REPORT_DIR)/covcounters go test ./...
	MOCKED=true go test -v -coverprofile=$(REPORT_DIR)/unit/acc_coverage.out -json ./... > $(REPORT_DIR)/unit/test_results.json;
	@echo "=== Generating reports ==="
	@if [ -f $(REPORT_DIR)/unit/acc_coverage.out ]; then \
		go tool cover -func=$(REPORT_DIR)/unit/acc_coverage.out > $(REPORT_DIR)/unit/coverage_summary.txt; \
		go tool cover -html=$(REPORT_DIR)/unit/acc_coverage.out -o $(REPORT_DIR)/unit/coverage.html; \
		echo "Coverage reports generated:"; \
		echo "  - Summary: $(REPORT_DIR)/unit/coverage_summary.txt"; \
		echo "  - HTML: $(REPORT_DIR)/unit/coverage.html"; \
	else \
		echo "No coverage data generated"; \
	fi
	@if [ -f $(REPORT_DIR)/unit/test_results.json ]; then \
		echo "Test results: $(REPORT_DIR)/unit/test_results.json"; \
	fi
	@echo "=== Test execution completed ==="

global-coverage: prepare-go prepare-reports
	@if [ -f $(REPORT_DIR)/unit/acc_coverage.out ]; then \
		go tool covdata merge -i=$(REPORT_DIR)/covcounters -o=$(REPORT_DIR)/covcounters/merged; \
		go tool covdata percent -i=$(REPORT_DIR)/covcounters/merged > $(REPORT_DIR)/acceptance/coverage_acceptance.txt; \
		go tool covdata textfmt -i=$(REPORT_DIR)/covcounters -o=$(REPORT_DIR)/acceptance/acceptance_merged.out; \
		go tool cover -func=$(REPORT_DIR)/acceptance/acceptance_merged.out > $(REPORT_DIR)/acceptance/coverage_func_summary.txt; \
		if [ "$(OS)" = "Darwin" ]; then \
			sed -i '' 's|.*terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' $(REPORT_DIR)/acceptance/coverage_func_summary.txt; \
		else \
			sed -i 's|.*terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' $(REPORT_DIR)/acceptance/coverage_func_summary.txt; \
		fi; \
		go tool cover -func=$(REPORT_DIR)/unit/acc_coverage.out > $(REPORT_DIR)/unit/coverage_summary.txt; \
		awk '/^total:/ {gsub(/%/, "", $$3); print $$3}' $(REPORT_DIR)/acceptance/coverage_func_summary.txt > $(REPORT_DIR)/acceptance/coverage_acceptance_summary.txt; \
		awk '/^total:/ {gsub(/%/, "", $$3); print $$3}' $(REPORT_DIR)/unit/coverage_summary.txt > $(REPORT_DIR)/unit/coverage_global.txt; \
		echo "Global coverage unit: $$(cat $(REPORT_DIR)/unit/coverage_global.txt)%"; \
		echo "Global coverage acceptance: $$(cat $(REPORT_DIR)/acceptance/coverage_acceptance_summary.txt)%"; \
	else \
		echo "0" > $(REPORT_DIR)/unit/coverage_global.txt; \
		echo "No coverage data found"; \
	fi

cobertura: prepare-go prepare-reports
	@if [ -f $(REPORT_DIR)/unit/acc_coverage.out ]; then \
		if [ "$(OS)" = "Darwin" ]; then \
			sed -i '' 's|.*terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' $(REPORT_DIR)/acceptance/acceptance_merged.out; \
		else \
			sed -i 's|.*terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' $(REPORT_DIR)/acceptance/acceptance_merged.out; \
		fi; \
		go tool github.com/boumenot/gocover-cobertura < $(REPORT_DIR)/unit/acc_coverage.out > $(REPORT_DIR)/unit/cobertura.xml; \
		go tool github.com/boumenot/gocover-cobertura < $(REPORT_DIR)/acceptance/acceptance_merged.out > $(REPORT_DIR)/acceptance/cobertura.xml; \
		echo "Cobertura report generated: $(REPORT_DIR)/cobertura.xml"; \
	else \
		echo "No coverage data found"; \
	fi

install:
	@go build -cover -o ${HOME}/go/bin/terraform-provider-infomaniak main.go

.PHONY: prepare-go prepare-logs prepare-reports testacc global-coverage cobertura install