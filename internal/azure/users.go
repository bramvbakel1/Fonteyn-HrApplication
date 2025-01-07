package azure

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents a user in Microsoft Graph API.
type User struct {
	ID                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	UserPrincipalName string   `json:"userPrincipalName"`
	FirstName         string   `json:"givenName"`
	LastName          string   `json:"surname"`
	Roles             []string `json:"roles"`
}

// FetchUsers retrieves a list of users from Microsoft Graph API.
func FetchUsers(accessToken string) ([]User, error) {
	url := "https://graph.microsoft.com/v1.0/users"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Value []User `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	// Fetch roles for each user
	for i := range result.Value {
		roles, err := fetchUserRoles(accessToken, result.Value[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching roles for user %s: %v", result.Value[i].ID, err)
		}
		result.Value[i].Roles = roles
	}

	return result.Value, nil
}

// fetchUserRoles retrieves roles assigned to a specific user from Microsoft Graph API.
func fetchUserRoles(accessToken, userID string) ([]string, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/memberOf", userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching roles for user %s: %v", userID, err)
	}
	defer resp.Body.Close()

	var result struct {
		Value []struct {
			DisplayName string `json:"displayName"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response body for roles: %v", err)
	}

	// Extract role names
	var roles []string
	for _, group := range result.Value {
		roles = append(roles, group.DisplayName)
	}

	return roles, nil
}

// DeleteUser deletes a user from Microsoft Graph API.
func DeleteUser(accessToken, userID string) error {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
