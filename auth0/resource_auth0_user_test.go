package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAuth0User(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0UserExists("auth0_user.test_user"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "connection_type", "Username-Password-Authentication"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "email", "test@example.com"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "name", "user1234"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "password", "8aabf4be-2ad5-48b6-84aa-3dcd112716f0"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "user_metadata.item1", "value1"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "user_metadata.item2", "value2"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "email_verified", "true"),
				),
			},
			{
				Config: testUpdateUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuth0UserExists("auth0_user.test_user"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "connection_type", "Username-Password-Authentication"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "email", "foo@example.com"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "name", "user1234"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "password", "8aabf4be-2ad5-48b6-84aa-3dcd112716f0"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "user_metadata.item1", "value4"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "user_metadata.item2", "value3"),
					resource.TestCheckResourceAttr("auth0_user.test_user", "email_verified", "true"),
				),
			},
		},
	})
}

func TestAccAuth0UserImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuth0UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testImportUserConfig,
			},
			{
				ResourceName:            "auth0_user.test_user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckAuth0UserDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*AuthClient)

	users := getResourcesByType("auth0_user", state)

	if len(users) != 1 {
		return fmt.Errorf("expecting only 1 auth0 user resource found %v", len(users))
	}

	response, err := client.GetUserById(users[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get auth0 user by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("user %s still exists, %+v", users[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckAuth0UserExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*AuthClient)

		user, err := client.GetUserById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if user == nil {
			return fmt.Errorf("user with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateUserConfig = `
resource "auth0_user" "test_user" {
	connection_type = "Username-Password-Authentication"
	email 			= "test@example.com"
	name 			= "user1234"
	password 		= "8aabf4be-2ad5-48b6-84aa-3dcd112716f0"
	user_metadata 	= {
		item1 = "value1"
		item2 = "value2"
	}
	email_verified 	= true
}
`

const testUpdateUserConfig = `
resource "auth0_user" "test_user" {
	connection_type = "Username-Password-Authentication"
	email 			= "foo@example.com"
	name 			= "user1234"
	password 		= "8aabf4be-2ad5-48b6-84aa-3dcd112716f0"
	user_metadata 	= {
		item1 = "value4"
		item2 = "value3"
	}
	email_verified 	= true
}
`

const testImportUserConfig = `

resource "auth0_user" "test_user" {
	connection_type = "Username-Password-Authentication"
	email 			= "foo@example.com"
	name 			= "user1234"
	password = "8aabf4be-2ad5-48b6-84aa-3dcd112716f0"
	user_metadata 	= {
		item1 = "value4"
		item2 = "value3"
	}
	email_verified 	= true
}
`
