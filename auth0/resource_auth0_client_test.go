package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAuth0Client(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0ClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateClientConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ClientExists("auth0_client.test_client"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "name", "test client"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "app_type", "non_interactive"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "grant_types.0", "password"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "token_endpoint_auth_method", "none"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "client_metadata.item1", "value1"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "client_metadata.item2", "value2"),
				),
			},
			{
				Config: testUpdateClientConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ClientExists("auth0_client.test_client"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "name", "test client"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "app_type", "non_interactive"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "grant_types.0", "password"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "token_endpoint_auth_method", "none"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "client_metadata.item1", "value3"),
					resource.TestCheckResourceAttr("auth0_client.test_client", "client_metadata.item2", "value4"),
				),
			},
		},
	})
}

func TestAccAuth0ClientImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0ClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testImportClientConfig,
			},
			{
				ResourceName:      "auth0_client.test_client",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAuth0ClientDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*AuthClient)

	auth0Clients := getResourcesByType("auth0_client", state)

	if len(auth0Clients) != 1 {
		return fmt.Errorf("expecting only 1 auth0 client resource found %v", len(auth0Clients))
	}

	response, err := client.GetClientById(auth0Clients[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get auth0 client by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("client %s still exists, %+v", auth0Clients[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckAuth0ClientExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*AuthClient)

		auth0Client, err := client.GetClientById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if auth0Client == nil {
			return fmt.Errorf("client with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateClientConfig = `
resource "auth0_client" "test_client" {
	name						= "test client"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value1"
		item2 = "value2"
	}
}
`

const testUpdateClientConfig = `

resource "auth0_client" "test_client" {
	name						= "test client"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value3"
		item2 = "value4"
	}
}

`

const testImportClientConfig = `

resource "auth0_client" "test_client" {
	name						= "test client import"
	app_type 					= "non_interactive"
	grant_types 				= ["password"]
	token_endpoint_auth_method 	= "none"
    client_metadata				= {
		item1 = "value3"
		item2 = "value4"
	}
}

`
