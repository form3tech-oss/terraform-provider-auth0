GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test testacc

travisbuild: deps default

test: fmtcheck
	go test -v . ./auth0

testacc: fmtcheck
	go test -v ./auth0 -run="TestAcc"

build: fmtcheck vet testacc
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-auth0 ~/.terraform.d/plugins/terraform-provider-auth0
	@echo "Build succeeded"

build-gox: deps fmtcheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-auth0" .

release:
	go get github.com/goreleaser/goreleaser; \
    scripts/release.sh; \

deps:
	go get -u golang.org/x/net/context; \
    go get -u github.com/mitchellh/gox; \

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

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile