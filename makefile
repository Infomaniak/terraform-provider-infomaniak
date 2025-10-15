LOG_DIRS = ./internal/services/kaas/tmp ./internal/services/dbaas/tmp

prepare-logs:
	@mkdir -p $(LOG_DIRS) 2>/dev/null || true
	@for dir in $(LOG_DIRS); do \
		touch $$dir/terraform.log; \
	done

testacc: prepare-logs
	TF_ACC=1 go test -v ./...

test: prepare-logs
	go test -v ./...

.PHONY: prepare-logs testacc test
