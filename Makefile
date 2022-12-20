GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test testacc

test: fmtcheck
	go test -v . ./auth0

testacc: fmtcheck
	go test -v ./auth0 -run="TestAcc"

build-only:
	@go build -o terraform-provider-auth0
	@mkdir -p "$HOME/.terraform.d/plugins/"
	@cp "terraform-provider-auth0" "$HOME/.terraform.d/plugins/terraform-provider-auth0"
	@echo "Build succeeded"

build: fmtcheck vet testacc build-only

clean:
	rm -rf pkg/

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build build-only test testacc vet fmt fmtcheck errcheck vendor-status test-compile
