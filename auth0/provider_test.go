package auth0

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	"os"
)


var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"auth0": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AUTH0_DOMAIN"); v == "" {
		t.Fatal("AUTH0_DOMAIN must be set for acceptance tests")
	}

	if v := os.Getenv("AUTH0_CLIENT_ID"); v == "" {
		t.Fatal("AUTH0_CLIENT_ID must be set for acceptance tests")
	}

	if v := os.Getenv("AUTH0_CLIENT_SECRET"); v == "" {
		t.Fatal("AUTH0_CLIENT_SECRET must be set for acceptance tests")
	}
}