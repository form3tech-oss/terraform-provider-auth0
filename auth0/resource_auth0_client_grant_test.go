package auth0

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAuth0ClientGrant(t *testing.T) {
	testUUID := uuid.New().String()

	testCreateClientGrantConfig := `
resource "auth0_client" "test_client" {
	name						= "test client grant ` + testUUID + `"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value1"
		item2 = "value2"
	}
}

resource "auth0_api" "test_api" {
	name 		= "https://api.example.com/client_grant_test` + testUUID + `"
	identifier 	= "https://api.example.com/client_grant_test` + testUUID + `"
}

resource "auth0_client_grant" "test_client_grant" {
	client_id 	= "${auth0_client.test_client.id}"
	audience 	= "https://api.example.com/client_grant_test` + testUUID + `"
	scope 		= ["something"]
}
`

	testCreateClientGrantConfigNoScope := `
resource "auth0_client" "test_client_no_scope" {
	name						= "test client grant no scope ` + testUUID + `"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value1"
		item2 = "value2"
	}
}

resource "auth0_api" "test_api_no_scope" {
	name 		= "https://api.example.com/client_grant_test_no_scope` + testUUID + `"
	identifier 	= "https://api.example.com/client_grant_test_no_scope` + testUUID + `"
}

resource "auth0_client_grant" "test_client_grant_no_scope" {
	client_id 	= "${auth0_client.test_client_no_scope.id}"
	audience 	= "https://api.example.com/client_grant_test` + testUUID + `"
}
`

	testUpdateClientGrantConfig := `

resource "auth0_client" "test_client" {
	name						= "test client ` + testUUID + `"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value1"
		item2 = "value2"
	}
}

resource "auth0_api" "test_api" {
	name 		= "https://api.example.com/client_grant_test` + testUUID + `"
	identifier 	= "https://api.example.com/client_grant_test` + testUUID + `"
}

resource "auth0_client_grant" "test_client_grant" {
	client_id 	= "${auth0_client.test_client.id}"
	audience 	= "https://api.example.com/client_grant_test` + testUUID + `"
	scope 		= ["something_else"]
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0ClientGrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateClientGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ClientGrantExists("auth0_client_grant.test_client_grant"),
					resource.TestCheckResourceAttrSet("auth0_client_grant.test_client_grant", "client_id"),
					resource.TestCheckResourceAttr("auth0_client_grant.test_client_grant", "audience", "https://api.example.com/client_grant_test"+testUUID),
					resource.TestCheckResourceAttr("auth0_client_grant.test_client_grant", "scope.0", "something"),
				),
			},
			{
				Config: testCreateClientGrantConfigNoScope,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ClientGrantExists("auth0_client_grant.test_client_grant_no_scope"),
					resource.TestCheckResourceAttrSet("auth0_client_grant.test_client_grant_no_scope", "client_id"),
					resource.TestCheckResourceAttr("auth0_client_grant.test_client_grant_no_scope", "audience", "https://api.example.com/client_grant_test_no_scope"+testUUID),
					resource.TestCheckNoResourceAttr("auth0_client_grant.test_client_grant_no_scope", "scope"),
				),
			},
			{
				Config: testUpdateClientGrantConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ClientGrantExists("auth0_client_grant.test_client_grant"),
					resource.TestCheckResourceAttrSet("auth0_client_grant.test_client_grant", "client_id"),
					resource.TestCheckResourceAttr("auth0_client_grant.test_client_grant", "audience", "https://api.example.com/client_grant_test"+testUUID),
					resource.TestCheckResourceAttr("auth0_client_grant.test_client_grant", "scope.0", "something_else"),
				),
			},
		},
	})
}

func testAccCheckAuth0ClientGrantDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*AuthClient)

	auth0ClientGrants := getResourcesByType("auth0_client_grant", state)

	if len(auth0ClientGrants) != 1 {
		return fmt.Errorf("expecting only 1 auth0 client-grant resource found %v", len(auth0ClientGrants))
	}

	clientId := auth0ClientGrants[0].Primary.Attributes["client_id"]
	audience := auth0ClientGrants[0].Primary.Attributes["audience"]

	auth0ClientGrant, err := client.GetClientGrantByClientIdAndAudience(clientId, audience)

	if err != nil {
		return fmt.Errorf("error calling get auth0 client-grant by id: %v", err)
	}

	if auth0ClientGrant != nil {
		return fmt.Errorf("client-grant %s still exists, %+v", auth0ClientGrants[0].Primary.ID, auth0ClientGrant)
	}

	return nil
}

func testAccCheckAuth0ClientGrantExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*AuthClient)

		clientId := rs.Primary.Attributes["client_id"]
		audience := rs.Primary.Attributes["audience"]

		auth0ClientGrant, err := client.GetClientGrantByClientIdAndAudience(clientId, audience)

		if err != nil {
			return err
		}

		if auth0ClientGrant == nil {
			return fmt.Errorf("client-grant with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}
