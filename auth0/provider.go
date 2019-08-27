package auth0

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/parnurzeal/gorequest"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_DOMAIN", nil),
			},
			"auth0_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
			},
			"auth0_client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"auth0_user":         resourceAuth0User(),
			"auth0_client":       resourceAuth0Client(),
			"auth0_api":          resourceAuth0Api(),
			"auth0_client_grant": resourceAuth0ClientGrant(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type Config struct {
	domain      string
	accessToken string
	apiUri      string
}

type LoginRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Auth0 client")

	domain := d.Get("domain").(string)
	apiUri := "https://" + domain + "/api/v2/"

	auth0LoginRequest := &LoginRequest{
		ClientId:     d.Get("auth0_client_id").(string),
		ClientSecret: d.Get("auth0_client_secret").(string),
		Audience:     apiUri,
		GrantType:    "client_credentials",
	}

	resp, body, errs := gorequest.New().Post("https://" + domain + "/oauth/token").Send(auth0LoginRequest).End()

	if errs != nil {
		return nil, fmt.Errorf("could not log in to auth0, error: %v", errs)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("auth0 /oauth/token request failed. status: %d, body: %s", resp.StatusCode, body)
	}

	loginResponse := &LoginResponse{}
	err := json.Unmarshal([]byte(body), loginResponse)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 login response, error: %v %s", err, body)
	}

	config := &Config{
		domain:      domain,
		accessToken: loginResponse.AccessToken,
		apiUri:      apiUri,
	}

	return NewClient(config), nil
}
