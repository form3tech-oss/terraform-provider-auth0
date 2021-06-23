package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuth0Client() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuth0ClientCreate,
		Read:   resourceAuth0ClientRead,
		Delete: resourceAuth0ClientDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"app_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"grant_types": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
				ForceNew: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAuth0ClientCreate(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	clientRequest := createClientRequestFromResourceData(d)

	client, err := auth0Client.CreateClient(clientRequest)

	if err != nil {
		return fmt.Errorf("failed to create auth0 client: %v error: %v", clientRequest, err)
	}

	d.SetId(client.ClientId)

	return resourceAuth0ClientRead(d, meta)
}

func resourceAuth0ClientRead(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	client, err := auth0Client.GetClientById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find auth0 client: %v", err)
	}

	if client == nil {
		d.SetId("")
	} else {
		d.Set("name", client.Name)
		d.Set("grant_types", client.GrantTypes)
		d.Set("app_type", client.ApplicationType)
		d.Set("token_endpoint_auth_method", client.TokenEndpointAuthMethod)
		d.Set("client_metadata", client.ClientMetaData)
		d.Set("client_secret", client.ClientSecret)
	}

	return nil
}

func resourceAuth0ClientDelete(d *schema.ResourceData, meta interface{}) error {

	auth0Client := meta.(*AuthClient)

	err := auth0Client.DeleteClientById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete auth0 client: %v", err)
	}

	return nil
}

func createClientRequestFromResourceData(d *schema.ResourceData) *ClientRequest {
	clientRequest := &ClientRequest{}

	clientRequest.Name = readStringFromResource(d, "name")
	clientRequest.ApplicationType = readStringFromResource(d, "app_type")
	clientRequest.GrantTypes = readStringArrayFromResource(d, "grant_types")
	clientRequest.TokenEndpointAuthMethod = readStringFromResource(d, "token_endpoint_auth_method")
	clientRequest.ClientMetaData = readMapFromResource(d, "client_metadata")

	return clientRequest
}
