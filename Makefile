default: build

build:
	go install

dev: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/gofireflyio/firefly/1.0.0/darwin_arm64
	cp $(GOPATH)/bin/terraform-provider-firefly ~/.terraform.d/plugins/registry.terraform.io/gofireflyio/firefly/1.0.0/darwin_arm64/

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 30m

# Run unit tests
test:
	go test ./... -timeout=10m

# Format code
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./

# Check formatting
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

# Vet code
vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

# Generate documentation
docs:
	@echo "==> Generating docs..."
	go generate

# Clean build artifacts
clean:
	rm -f terraform-provider-firefly

# Install development dependencies
deps:
	go mod download
	go mod tidy

# Quick development cycle
devbuild: fmt vet build

# Debug build
debug: 
	go build -gcflags="all=-N -l" -o terraform-provider-firefly

.PHONY: build dev testacc test fmt fmtcheck vet docs clean deps devbuild debug