LOG_DIRS = ./internal/services/kaas/tmp ./internal/services/dbaas/tmp
REPORT_DIR = reports

prepare-go:
	@go mod download

prepare-logs:
	@mkdir -p $(LOG_DIRS) 2>/dev/null || true
	@for dir in $(LOG_DIRS); do \
		touch $$dir/terraform.log; \
	done

prepare-reports:
	@mkdir -p $(REPORT_DIR)

testacc: prepare-go prepare-logs prepare-reports
	@echo "=== Running acceptance tests ==="
	go test -v -coverprofile=$(REPORT_DIR)/acc_coverage.out -json ./... > $(REPORT_DIR)/test_results.json;
	@echo "=== Generating reports ==="
	@if [ -f $(REPORT_DIR)/acc_coverage.out ]; then \
		go tool cover -func=$(REPORT_DIR)/acc_coverage.out > $(REPORT_DIR)/coverage_summary.txt; \
		go tool cover -html=$(REPORT_DIR)/acc_coverage.out -o $(REPORT_DIR)/coverage.html; \
		echo "Coverage reports generated:"; \
		echo "  - Summary: $(REPORT_DIR)/coverage_summary.txt"; \
		echo "  - HTML: $(REPORT_DIR)/coverage.html"; \
	else \
		echo "No coverage data generated"; \
	fi
	@if [ -f $(REPORT_DIR)/test_results.json ]; then \
		echo "Test results: $(REPORT_DIR)/test_results.json"; \
	fi
	@echo "=== Test execution completed ==="

global-coverage: prepare-go prepare-reports
	@if [ -f $(REPORT_DIR)/acc_coverage.out ]; then \
		go tool cover -func=$(REPORT_DIR)/acc_coverage.out > $(REPORT_DIR)/coverage_summary.txt; \
		awk '/^total:/ {gsub(/%/, "", $$3); print $$3}' $(REPORT_DIR)/coverage_summary.txt > $(REPORT_DIR)/coverage_global.txt; \
		echo "Global coverage: $$(cat $(REPORT_DIR)/coverage_global.txt)%"; \
	else \
		echo "0" > $(REPORT_DIR)/coverage_global.txt; \
		echo "No coverage data found"; \
	fi

cobertura: prepare-go prepare-reports
	@if [ -f $(REPORT_DIR)/acc_coverage.out ]; then \
		go tool github.com/boumenot/gocover-cobertura < $(REPORT_DIR)/acc_coverage.out > $(REPORT_DIR)/cobertura.xml \
		echo "Cobertura report generated: $(REPORT_DIR)/cobertura.xml"; \
	else \
		echo "No coverage data found"; \
	fi

.PHONY: prepare-go prepare-logs prepare-reports testacc global-coverage cobertura