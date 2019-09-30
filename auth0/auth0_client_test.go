package auth0

import (
	"github.com/google/uuid"
	"os"
	"sync"
	"testing"
	"time"
)

// 1. Create Token
// 2. Create new user
// 3. Simulate throttled load to GetUserById
// 4. Clean up the created user
func TestAccGetUserByIdIsNotThrottled(t *testing.T) {
	auth0RetryCount := 2
	timeBeetwenRetries := time.Second
	numberOfRequests := 100
	numberOfGoRoutines := 10

	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		t.Fatal("AUTH0_DOMAIN must be set for acceptance tests")
	}

	clientId := os.Getenv("AUTH0_CLIENT_ID")
	if clientId == "" {
		t.Fatal("AUTH0_CLIENT_ID must be set for acceptance tests")
	}

	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	if clientSecret == "" {
		t.Fatal("AUTH0_CLIENT_SECRET must be set for acceptance tests")
	}

	apiUri := "https://" + domain + "/api/v2/"

	config := &Config{
		domain:             domain,
		apiUri:             apiUri,
		maxRetryCount:      auth0RetryCount,
		timeBetweenRetries: timeBeetwenRetries,
	}

	client, err := NewClient(clientId, clientSecret, config)
	if err != nil {
		t.Fatalf("auth0 test cliend creation failure %v", err)
	}

	userRequest := &UserRequest{
		Connection:    "Username-Password-Authentication",
		Email:         "auth0-provider-test@auth0-provider-test.com",
		Name:          "auth0-provider-test",
		Password:      uuid.New().String(),
		UserMetaData:  nil,
		EmailVerified: false,
	}

	createdUser, err := client.CreateUser(userRequest)

	if err != nil {
		t.Fatalf("failed to create test user %v", err)
	}

	defer func() {
		err := client.DeleteUserById(createdUser.UserId)
		if err != nil {
			t.Fatalf("Dangling resource! Failed to remove test user with UserId '%v' with error message: %v", createdUser.UserId, err)
		}
	}()

	var done sync.WaitGroup

	for i := 0; i < numberOfGoRoutines; i++ {
		done.Add(1)
		go func() {
			defer done.Done()
			for i := 1; i <= numberOfRequests; i++ {
				_, err := client.GetUserById(createdUser.UserId)
				if err != nil {
					t.Fatalf("failed to get user %v", err)
				}
			}
		}()
	}

	done.Wait()
}
