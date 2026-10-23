.PHONY: lint

lint:
	golangci-lint run -c .golangci.yml ./...
	@echo gopls check -severity hint ./...
	@RESULT=$$(grep -rL '^// Code generated .* DO NOT EDIT' --include="*.go" . | bash -c ' xargs gopls check -severity hint 2>&1'); \
	if [ -n "$$RESULT" ]; then \
		echo "$$RESULT"; \
		exit 1; \
	fi
