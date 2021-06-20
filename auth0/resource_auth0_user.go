package auth0

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuth0User() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuth0UserCreate,
		Read:   resourceAuth0UserRead,
		Update: resourceAuth0UserUpdate,
		Delete: resourceAuth0UserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
			},
			"email_verified": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAuth0UserCreate(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	userRequest := createUserRequestFromResourceData(d)

	user, err := auth0Client.CreateUser(userRequest)
	TfLogJson("[resourceAuth0UserCreate]", userRequest)

	if err != nil {
		return fmt.Errorf("failed to create auth0 user: %v error: %v", userRequest, err)
	}

	d.SetId(user.UserId)

	return resourceAuth0UserRead(d, meta)
}

func resourceAuth0UserUpdate(d *schema.ResourceData, meta interface{}) error {
	auth0Client := meta.(*AuthClient)

	updateUserRequests := createUserUpdatesFromResourceData(d)
	userId := d.Id()

	for _, update := range updateUserRequests {
		TfLogJson("[resourceAuth0UserUpdate]", update)
		_, err := auth0Client.UpdateUserById(userId, update)

		if err != nil {
			return err
		}
	}

	return nil
}

func resourceAuth0UserRead(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	user, err := auth0Client.GetUserById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find auth0 user: %v", err)
	}

	if user == nil {
		d.SetId("")
		TfLogString("[resourceAuth0UserRead]", "User is nil")
	} else {
		d.Set("user_id", user.UserId)
		d.Set("email", user.Email)
		d.Set("name", user.Name)
		d.Set("user_metadata", user.UserMetaData)
		d.Set("email_verified", user.EmailVerified)

		// TODO: We should model identities properly as more than one identity for a user
		// might exist
		if len(user.Identities) > 0 {
			d.Set("connection_type", user.Identities[0].Connection)
		}

		TfLogJson("[resourceAuth0UserRead]", user)
	}

	return nil
}

func resourceAuth0UserDelete(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	err := auth0Client.DeleteUserById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete auth0 user: %v", err)
	}

	TfLogString("[resourceAuth0UserDelete]", d.Id())

	return nil
}

// Given ResourceData from terraform, generate a list of patches to apply sequentially.
// We do this to reconcile state against the auth0 API, you can't update email_verified, and
// passsword, or email in one API call for example, so it must be updated in multiple requests.
func createUserUpdatesFromResourceData(d *schema.ResourceData) []*UserRequest {

	// The first patch contains everything but email_verified, and password state
	userRequestA := &UserRequest{}

	userRequestA.Connection = readStringFromResource(d, "connection_type")
	userRequestA.Email = readStringFromResource(d, "email")
	userRequestA.Name = readStringFromResource(d, "name")
	userRequestA.UserMetaData = readMapFromResource(d, "user_metadata")

	TfLogJson("[createUserUpdatesFromResourceData-userRequestA]", userRequestA)

	// Second contains only the password
	userRequestB := &UserRequest{}
	userRequestB.Password = readStringFromResource(d, "password")

	TfLogJson("[createUserUpdatesFromResourceData-userRequestB]", userRequestB)

	// Final updates the email_verified state
	userRequestC := &UserRequest{}
	userRequestC.EmailVerified = readBoolFromResource(d, "email_verified")

	TfLogJson("[createUserUpdatesFromResourceData-userRequestC]", userRequestC)

	return []*UserRequest{userRequestA, userRequestB, userRequestC}
}

func createUserRequestFromResourceData(d *schema.ResourceData) *UserRequest {

	userRequest := &UserRequest{}

	userRequest.Connection = readStringFromResource(d, "connection_type")
	userRequest.Email = readStringFromResource(d, "email")
	userRequest.Name = readStringFromResource(d, "name")
	userRequest.Password = readStringFromResource(d, "password")
	userRequest.UserMetaData = readMapFromResource(d, "user_metadata")
	userRequest.EmailVerified = readBoolFromResource(d, "email_verified")

	TfLogJson("[createUserRequestFromResourceData-userRequest]", userRequest)

	return userRequest
}
