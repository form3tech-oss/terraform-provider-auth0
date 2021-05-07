package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuth0Api() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuth0ApiCreate,
		Read:   resourceAuth0ApiRead,
		Delete: resourceAuth0ApiDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAuth0ApiCreate(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	apiRequest := createApiRequestFromResourceData(d)

	api, err := auth0Client.CreateApi(apiRequest)

	if err != nil {
		return fmt.Errorf("failed to create auth0 api: %v error: %v", apiRequest, err)
	}

	d.SetId(api.Id)

	return resourceAuth0ApiRead(d, meta)
}

func resourceAuth0ApiRead(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	api, err := auth0Client.GetApiById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find auth0 api: %v", err)
	}

	if api == nil {
		d.SetId("")
	} else {
		d.Set("name", api.Name)
		d.Set("identifier", api.Identifier)
	}

	return nil
}

func resourceAuth0ApiDelete(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	err := auth0Client.DeleteApiById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete auth0 api: %v", err)
	}

	return nil
}

func createApiRequestFromResourceData(d *schema.ResourceData) *ApiRequest {
	apiRequest := &ApiRequest{}

	apiRequest.Name = readStringFromResource(d, "name")
	apiRequest.Identifier = readStringFromResource(d, "identifier")

	return apiRequest
}
