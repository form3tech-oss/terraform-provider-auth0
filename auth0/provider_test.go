package auth0

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"auth0": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *Provider()
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
