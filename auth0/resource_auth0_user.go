package auth0

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func resourceAuth0User() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuth0UserCreate,
		Read:   resourceAuth0UserRead,
		Delete: resourceAuth0UserDelete,

		Schema: map[string]*schema.Schema{
			"connection_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
				ForceNew: true,
			},
			"email_verified": &schema.Schema{
				Type:         schema.TypeBool,
				Optional:     true,
				ForceNew: 	  true,
			},
		},
	}
}


func resourceAuth0UserCreate(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	userRequest := createUserRequestFromResourceData(d)

	user, err := auth0Client.CreateUser(userRequest)

	if err != nil {
		return fmt.Errorf("failed to create auth0 user: %v error: %v", userRequest, err)
	}

	d.SetId(user.UserId)

	return resourceAuth0UserRead(d, meta)
}

func resourceAuth0UserRead(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	user, err := auth0Client.GetUserById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find auth0 user: %v", err)
	}

	if user == nil {
		d.SetId("")
	} else {
		d.Set("email", user.Email)
		d.Set("name", user.Name)
		d.Set("user_metadata", user.UserMetaData)
		d.Set("email_verified", user.EmailVerified)
	}

	return nil
}

func resourceAuth0UserDelete(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	err := auth0Client.DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete auth0 user: %v", err)
	}

	return nil
}

func createUserRequestFromResourceData(d *schema.ResourceData) *UserRequest {

	userRequest := &UserRequest{}

	userRequest.Connection = readStringFromResource(d, "connection_type")
	userRequest.Email = readStringFromResource(d, "email")
	userRequest.Name = readStringFromResource(d, "name")
	userRequest.Password = readStringFromResource(d, "password")
	userRequest.UserMetaData = readMapFromResource(d, "user_metadata")
	userRequest.EmailVerified = readBoolFromResource(d, "email_verified")


	return userRequest
}