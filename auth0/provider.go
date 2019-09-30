package auth0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"time"
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
			"auth0_request_max_retry_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Max retry on requests to Auth0",
			},
			"auth0_time_between_retries": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//second, cannot use time time.Second here as under the hood provider framework users cty.Value dynamic types
				//time.Second is translated to 1000000000 (Nanoseconds) which ends up with error panic: can't convert 1000000000 to cty.Value
				Default:     1000,
				Description: "Time to wait between retried requests to Auth0 (in milliseconds)",
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
	domain             string
	accessToken        string
	apiUri             string
	maxRetryCount      int
	timeBetweenRetries time.Duration
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
	clientId := d.Get("auth0_client_id").(string)
	clientSecret := d.Get("auth0_client_secret").(string)
	maxRetryCount := d.Get("auth0_request_max_retry_count").(int)
	timeBetweenRetries := d.Get("auth0_request_max_retry_count").(int)

	config := &Config{
		domain:             domain,
		apiUri:             apiUri,
		maxRetryCount:      maxRetryCount,
		timeBetweenRetries: time.Duration(timeBetweenRetries) * time.Millisecond,
	}

	client, err := NewClient(clientId, clientSecret, config)

	if err != nil {
		return nil, fmt.Errorf("auth0 provider configuration failure, error: %v", err)
	}

	return client, nil
}
