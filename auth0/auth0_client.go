package auth0

import (
	"github.com/parnurzeal/gorequest"
	"fmt"
	"encoding/json"
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
	UserId			 string		`json:"user_id,omitempty"`
	Email 			 string     `json:"email,omitempty"`
	Name 		 	 string     `json:"name,omitempty"`
	UserMetaData 	 map[string]interface{}     `json:"user_metadata,omitempty"`
	EmailVerified 	 bool       `json:"email_verified,omitempty"`
}


func (config *Config) getAuthenticationHeader() string {
	return "Bearer " + config.accessToken
}

func (authClient *AuthClient) GetUserById(id string) (*User, error) {

	_, body, errs := gorequest.New().Get(authClient.config.apiUri + "users/" + id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()

	if errs != nil {
		return nil, fmt.Errorf("could parse user resposne from auth0, error: %v", errs)
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

	_, body, errs := gorequest.New().Post(authClient.config.apiUri + "users").Send(userRequest).Set("Authorization", authClient.config.getAuthenticationHeader()).End()


	if errs != nil {
		return nil, fmt.Errorf("could create user in auth0, error: %v", errs)
	}

	createdUser := &User{}
	err := json.Unmarshal([]byte(body), createdUser)
	if err != nil {
		return nil, fmt.Errorf("could not parse auth0 user creation response, error: %v %s", err, body)
	}

	if createdUser.UserId == "" {
		return nil, fmt.Errorf("could not create user, erorr: %s", body)
	}

	return createdUser, nil
}

func (authClient *AuthClient) UpdateUserById(id string, userRequest *UserRequest) (*User, error) {

	_, body, errs := gorequest.New().Patch(authClient.config.apiUri + "users/" + id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
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

func (authClient *AuthClient) DeleteById(id string) error {

	res, _, errs := gorequest.New().Delete(authClient.config.apiUri + "users/" + id).Set("Authorization", authClient.config.getAuthenticationHeader()).End()
	if errs != nil {
		return fmt.Errorf("could not delete auth0 user, result: %v error: %v", res, errs)
	}

	return nil
}





