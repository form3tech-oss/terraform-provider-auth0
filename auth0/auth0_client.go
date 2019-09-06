package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type AuthClient struct {
	config *Config
}

func NewClient(config *Config) *AuthClient {
	return &AuthClient{
		config: config,
	}
}

type UserRequest struct {
	Connection    string                 `json:"connection,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Password      string                 `json:"password,omitempty"`
	UserMetaData  map[string]interface{} `json:"user_metadata,omitempty"`
	EmailVerified bool                   `json:"email_verified,omitempty"`
}

type User struct {
	UserId        string                 `json:"user_id,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Name          string                 `json:"name,omitempty"`
	UserMetaData  map[string]interface{} `json:"user_metadata,omitempty"`
	EmailVerified bool                   `json:"email_verified,omitempty"`
	Identities    []Identity             `json:"identities,omitempty"`
}

type Identity struct {
	Connection string `json:"connection,omitempty"`
	UserId     string `json:"user_id,omitempty"`
	Provider   string `json:"provider,omitempty"`
	IsSocial   bool   `json:"isSocial,omitempty"`
}

type ClientRequest struct {
	Name                    string                 `json:"name,omitempty"`
	ApplicationType         string                 `json:"app_type,omitempty"`
	GrantTypes              []string               `json:"grant_types,omitempty"`
	TokenEndpointAuthMethod string                 `json:"token_endpoint_auth_method,omitempty"`
	ClientMetaData          map[string]interface{} `json:"client_metadata,omitempty"`
}

type Client struct {
	ClientId                string                 `json:"client_id,omitempty"`
	ClientSecret            string                 `json:"client_secret,omitempty"`
	Name                    string                 `json:"name,omitempty"`
	ApplicationType         string                 `json:"app_type,omitempty"`
	GrantTypes              []string               `json:"grant_types,omitempty"`
	TokenEndpointAuthMethod string                 `json:"token_endpoint_auth_method,omitempty"`
	ClientMetaData          map[string]interface{} `json:"client_metadata,omitempty"`
}

type ApiRequest struct {
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type Api struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type ClientGrantRequest struct {
	ClientId string   `json:"client_id,omitempty"`
	Audience string   `json:"audience,omitempty"`
	Scope    []string `json:"scope,omitempty"`
}

func (cgr *ClientGrantRequest) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString("{")

	b.WriteString(`"client_id": "` + cgr.ClientId + `",`)
	b.WriteString(`"audience": "` + cgr.Audience + `",`)
	b.WriteString(`"scope": [`)

	if cgr.Scope != nil {
		for idx, scope := range cgr.Scope {
			if idx > 0 {
				b.WriteRune(',')
			}

			b.WriteString(`"` + scope + `"`)
		}
	}

	b.WriteString("]}")
	return b.Bytes(), nil
}

type ClientGrant struct {
	Id       string   `json:"id,omitempty"`
	ClientId string   `json:"client_id,omitempty"`
	Audience string   `json:"audience,omitempty"`
	Scope    []string `json:"scope,omitempty"`
}

func (config *Config) getAuthenticationHeader() string {
	return "Bearer " + config.accessToken
}

// User
func (authClient *AuthClient) GetUserById(id string) (*User, error) {

	resp, body, errs := gorequest.New().Get(authClient.config.apiUri+"users/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could parse user response from auth0, error: %v", errs)
	}

	user := &User{}
	err := json.Unmarshal([]byte(body), user)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 get user response, error: %v %s", err, body)
	}

	if user.UserId == "" {
		return nil, nil
	}

	return user, nil
}

func (authClient *AuthClient) CreateUser(userRequest *UserRequest) (*User, error) {

	resp, body, errs := gorequest.New().Post(authClient.config.apiUri+"users").Send(userRequest).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could create user in auth0, error: %v", errs)
	}

	createdUser := &User{}
	err := json.Unmarshal([]byte(body), createdUser)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 user creation response, error: %v %s", err, body)
	}

	if createdUser.UserId == "" {
		return nil, fmt.Errorf("could not create user, error: %s", body)
	}

	return createdUser, nil
}

func (authClient *AuthClient) UpdateUserById(id string, userRequest *UserRequest) (*User, error) {

	resp, body, errs := gorequest.New().
		Patch(authClient.config.apiUri+"users/"+id).
		Set("Authorization", authClient.config.getAuthenticationHeader()).
		Set("Content-Type", "application/json").
		Send(userRequest).
		End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could not update auth0 user, error: %v", errs)
	}

	updatedUser := &User{}
	err := json.Unmarshal([]byte(body), updatedUser)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 user update response, error: %v", err)
	}

	if updatedUser.UserId == "" {
		return nil, fmt.Errorf("could not update auth0 user, error: %v", body)
	}

	return updatedUser, nil
}

func (authClient *AuthClient) DeleteUserById(id string) error {

	res, _, errs := gorequest.New().Delete(authClient.config.apiUri+"users/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
	if errs != nil {
		return fmt.Errorf("could not delete auth0 user, result: %v error: %v", res, errs)
	}

	return nil
}

// Client
func (authClient *AuthClient) GetClientById(id string) (*Client, error) {

	resp, body, errs := gorequest.New().Get(authClient.config.apiUri+"clients/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could parse client response from auth0, error: %v", errs)
	}

	client := &Client{}
	err := json.Unmarshal([]byte(body), client)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 get client response, error: %v %s", err, body)
	}

	if client.ClientId == "" {
		return nil, nil
	}

	return client, nil
}

func (authClient *AuthClient) CreateClient(clientRequest *ClientRequest) (*Client, error) {

	resp, body, errs := gorequest.New().Post(authClient.config.apiUri+"clients").Send(clientRequest).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could create client in auth0, error: %v", errs)
	}

	createdClient := &Client{}
	err := json.Unmarshal([]byte(body), createdClient)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 client creation response, error: %v %s", err, body)
	}

	if createdClient.ClientId == "" {
		return nil, fmt.Errorf("could not create client, error: %s", body)
	}

	return createdClient, nil
}

func (authClient *AuthClient) UpdateClientById(id string, clientRequest *ClientRequest) (*Client, error) {
	resp, body, errs := gorequest.New().Patch(authClient.config.apiUri+"clients/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could not update auth0 client, error: %v", errs)
	}

	updatedClient := &Client{}
	err := json.Unmarshal([]byte(body), updatedClient)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 client update response, error: %v", err)
	}

	if updatedClient.ClientId == "" {
		return nil, fmt.Errorf("could not update auth0 client, error: %v", body)
	}

	return updatedClient, nil
}

func (authClient *AuthClient) DeleteClientById(id string) error {
	res, _, errs := gorequest.New().Delete(authClient.config.apiUri+"clients/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
	if errs != nil {
		return fmt.Errorf("could not delete auth0 client, result: %v error: %v", res, errs)
	}

	return nil
}

// Api
func (authClient *AuthClient) GetApiById(id string) (*Api, error) {

	resp, body, errs := gorequest.New().Get(authClient.config.apiUri+"resource-servers/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could parse api response from auth0, error: %v", errs)
	}

	api := &Api{}
	err := json.Unmarshal([]byte(body), api)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 get api response, error: %v %s", err, body)
	}

	if api.Id == "" {
		return nil, nil
	}

	return api, nil
}

func (authClient *AuthClient) CreateApi(apiRequest *ApiRequest) (*Api, error) {

	resp, body, errs := gorequest.New().Post(authClient.config.apiUri+"resource-servers").Send(apiRequest).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could create api in auth0, error: %v", errs)
	}

	createdApi := &Api{}
	err := json.Unmarshal([]byte(body), createdApi)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 api creation response, error: %v %s", err, body)
	}

	if createdApi.Id == "" {
		return nil, fmt.Errorf("could not create api, error: %s", body)
	}

	return createdApi, nil
}

func (authClient *AuthClient) UpdateApiById(id string, apiRequest *ApiRequest) (*Api, error) {

	resp, body, errs := gorequest.New().Patch(authClient.config.apiUri+"resource-servers/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could not update auth0 api, error: %v", errs)
	}

	updatedApi := &Api{}
	err := json.Unmarshal([]byte(body), updatedApi)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 api update response, error: %v", err)
	}

	if updatedApi.Id == "" {
		return nil, fmt.Errorf("could not update auth0 api, error: %v", body)
	}

	return updatedApi, nil
}

func (authClient *AuthClient) DeleteApiById(id string) error {

	res, _, errs := gorequest.New().Delete(authClient.config.apiUri+"resource-servers/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
	if errs != nil {
		return fmt.Errorf("could not delete auth0 api, result: %v error: %v", res, errs)
	}

	return nil
}

// GetClientGrantById retrieves a client grant based on ID.
//
// Note that this is significantly heavier than GetClientGrantByClientIdAndAudience due to the Auth0 API's lack of
// ID-based retrieval.
func (authClient *AuthClient) GetClientGrantById(id string) (*ClientGrant, error) {

	_, body, errs := gorequest.New().
		Get(authClient.config.apiUri+"client-grants").
		Set("Authorization", authClient.config.getAuthenticationHeader()).
		End()

	if errs != nil {
		return nil, fmt.Errorf("could parse client-grant response from auth0, error: %v", errs)
	}

	clientGrant := make([]ClientGrant, 0)
	err := json.Unmarshal([]byte(body), &clientGrant)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 get client-grant response, error: %v %s", err, body)
	}

	for _, grant := range clientGrant {
		if grant.Id == id {
			return &grant, nil
		}
	}

	return nil, nil
}

// ClientGrant
func (authClient *AuthClient) GetClientGrantByClientIdAndAudience(clientId string, audience string) (*ClientGrant, error) {

	queryParams := map[string]string{
		"client_id": clientId,
		"audience":  audience,
	}

	resp, body, errs := gorequest.New().
		Get(authClient.config.apiUri+"client-grants").
		Query(queryParams).Set("Authorization", authClient.config.getAuthenticationHeader()).
		End()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could parse client-grant response from auth0, error: %v", errs)
	}

	clientGrant := make([]ClientGrant, 0)
	err := json.Unmarshal([]byte(body), &clientGrant)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 get client-grant response, error: %v %s", err, body)
	}

	if len(clientGrant) != 1 || clientGrant[0].Id == "" {
		return nil, nil
	}

	return &clientGrant[0], nil
}

func (authClient *AuthClient) CreateClientGrant(clientGrantRequest *ClientGrantRequest) (*ClientGrant, error) {
	reqJSON, err := json.Marshal(clientGrantRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client grant request: %v", err)
	}

	resp, body, errs := gorequest.New().Post(authClient.config.apiUri+"client-grants").SendString(string(reqJSON)).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could create client-grant in auth0, error: %v", errs)
	}

	createdClientGrant := &ClientGrant{}
	err = json.Unmarshal([]byte(body), createdClientGrant)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 client-grant creation response, error: %v %s", err, body)
	}

	if createdClientGrant.Id == "" {
		return nil, fmt.Errorf("could not create client-grant, error: %s", body)
	}

	return createdClientGrant, nil
}

func (authClient *AuthClient) UpdateClientGrantById(id string, clientGrantRequest *ClientGrantRequest) (*ClientGrant, error) {

	resp, body, errs := gorequest.New().Patch(authClient.config.apiUri+"client-grants/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code (%d): %s", resp.StatusCode, body)
	}

	if errs != nil {
		return nil, fmt.Errorf("could not update auth0 client-grant, error: %v", errs)
	}

	updatedClientGrant := &ClientGrant{}
	err := json.Unmarshal([]byte(body), updatedClientGrant)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 client-grant update response, error: %v", err)
	}

	if updatedClientGrant.Id == "" {
		return nil, fmt.Errorf("could not update auth0 client-grant, error: %v", body)
	}

	return updatedClientGrant, nil
}

func (authClient *AuthClient) DeleteClientGrantById(id string) error {

	res, _, errs := gorequest.New().Delete(authClient.config.apiUri+"client-grants/"+id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
	if errs != nil {
		return fmt.Errorf("could not delete auth0 client-grant, result: %v error: %v", res, errs)
	}

	return nil
}
