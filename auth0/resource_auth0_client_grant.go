package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuth0ClientGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuth0ClientGrantCreate,
		Read:   resourceAuth0ClientGrantRead,
		Delete: resourceAuth0ClientGrantDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAuth0ClientGrantCreate(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	clientGrantRequest := createClientGrantRequestFromResourceData(d)

	clientGrant, err := auth0Client.CreateClientGrant(clientGrantRequest)

	if err != nil {
		return fmt.Errorf("failed to create auth0 client-grant: %v error: %v", clientGrantRequest, err)
	}

	fmt.Println("Created grant with id", clientGrant.Id)

	d.SetId(clientGrant.Id)

	return resourceAuth0ClientGrantRead(d, meta)
}

func resourceAuth0ClientGrantRead(d *schema.ResourceData, meta interface{}) error {
	var clientGrant *ClientGrant
	var err error

	auth0Client := meta.(*AuthClient)

	clientId := readStringFromResource(d, "client_id")
	audience := readStringFromResource(d, "audience")

	if clientId != "" && audience != "" {
		clientGrant, err = auth0Client.GetClientGrantByClientIdAndAudience(clientId, audience)
	} else {
		// This is necessary for ID-only import but it's significantly heavier than querying by client ID and audience
		clientGrant, err = auth0Client.GetClientGrantById(d.Id())
	}

	if err != nil {
		return fmt.Errorf("could not find auth0 client-grant: %v", err)
	}

	//fmt.Println("Read after create", clientGrant)

	if clientGrant == nil {
		d.SetId("")
	} else {
		d.Set("client_id", clientGrant.ClientId)
		d.Set("audience", clientGrant.Audience)
		d.Set("scope", clientGrant.Scope)
	}

	return nil
}

func resourceAuth0ClientGrantDelete(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	err := auth0Client.DeleteClientGrantById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete auth0 client-grant: %v", err)
	}

	return nil
}

func createClientGrantRequestFromResourceData(d *schema.ResourceData) *ClientGrantRequest {
	clientGrantRequest := &ClientGrantRequest{}

	clientGrantRequest.ClientId = readStringFromResource(d, "client_id")
	clientGrantRequest.Audience = readStringFromResource(d, "audience")
	clientGrantRequest.Scope = readStringArrayFromResource(d, "scope")

	return clientGrantRequest
}
