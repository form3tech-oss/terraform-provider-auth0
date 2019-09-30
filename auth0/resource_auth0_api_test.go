package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccAuth0Api(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0ApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ApiExists("auth0_api.test_api"),
					resource.TestCheckResourceAttr("auth0_api.test_api", "name", "https://api.example.com/api_test/1"),
					resource.TestCheckResourceAttr("auth0_api.test_api", "identifier", "https://api.example.com/api_test"),
				),
			},
			{
				Config: testUpdateApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0ApiExists("auth0_api.test_api"),
					resource.TestCheckResourceAttr("auth0_api.test_api", "name", "https://api.example.com/api_test/2"),
					resource.TestCheckResourceAttr("auth0_api.test_api", "identifier", "https://api.example.com/api_test"),
				),
			},
		},
	})
}

func TestAccAuth0ApiImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0ApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testImportApiConfig,
			},
			{
				ResourceName:      "auth0_api.test_api",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAuth0ApiDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*AuthClient)

	auth0Apis := getResourcesByType("auth0_api", state)

	if len(auth0Apis) != 1 {
		return fmt.Errorf("expecting only 1 auth0 api resource found %v", len(auth0Apis))
	}

	response, err := client.GetApiById(auth0Apis[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get auth0 api by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("api %s still exists, %+v", auth0Apis[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckAuth0ApiExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*AuthClient)

		auth0Api, err := client.GetApiById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if auth0Api == nil {
			return fmt.Errorf("api with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateApiConfig = `

resource "auth0_api" "test_api" {
	name 		= "https://api.example.com/api_test/1"
	identifier 	= "https://api.example.com/api_test"
}

`

const testUpdateApiConfig = `

resource "auth0_api" "test_api" {
	name 		= "https://api.example.com/api_test/2"
	identifier 	= "https://api.example.com/api_test"
}

`

const testImportApiConfig = `

resource "auth0_api" "test_api" {
	name 		= "https://api.example.com/api_test/2"
	identifier 	= "https://api.example.com/api_test_import"
}

`
