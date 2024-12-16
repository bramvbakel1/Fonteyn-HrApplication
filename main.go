package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// User struct represents a user in Microsoft Graph API
type User struct {
	ID                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	UserPrincipalName string   `json:"userPrincipalName"`
	FirstName         string   `json:"givenName"`
	LastName          string   `json:"surname"`
	Roles             []string `json:"roles"` // Store roles as a slice of strings
}

// Load environment variables and return them
func loadEnv() (string, string, string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("error loading .env file")
	}

	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	graphAPIURL := os.Getenv("GRAPH_API_URL")
	certPath := os.Getenv("SSL_CERT_PATH")
	keyPath := os.Getenv("SSL_KEY_PATH")

	return tenantID, clientID, clientSecret, graphAPIURL, certPath, keyPath, nil
}

// GetAccessToken retrieves the access token from Microsoft identity platform using OAuth2 client credentials flow
func GetAccessToken(tenantID, clientID, clientSecret string) (string, error) {
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token",
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	token, err := conf.Token(oauth2.NoContext)
	if err != nil {
		return "", fmt.Errorf("unable to get token: %v", err)
	}

	return token.AccessToken, nil
}

// Serve the homepage
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/home.html")
	if err != nil {
		http.Error(w, "Error loading homepage", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Serve the users page
func serveUsersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/users.html")
	if err != nil {
		http.Error(w, "Error loading users page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Handle fetching users from Microsoft Graph API
func handleUsers(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	tenantID, clientID, clientSecret, _, _, _, err := loadEnv()
	if err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		return
	}

	// Fetch access token
	accessToken, err := GetAccessToken(tenantID, clientID, clientSecret)
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	// Make API request to fetch users
	users, err := fetchUsers(accessToken)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	// Return users as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Fetch users from Microsoft Graph API
func fetchUsers(accessToken string) ([]User, error) {
	url := "https://graph.microsoft.com/v1.0/users"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Fetch the users
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

	for i := range result.Value {
		roles, err := fetchUserRoles(accessToken, result.Value[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching roles for user %s: %v", result.Value[i].ID, err)
		}
		result.Value[i].Roles = roles
	}

	return result.Value, nil
}

// Fetch roles assigned to a specific user
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

// Main server initialization
func main() {
	// Load environment variables
	_, _, _, _, certPath, keyPath, err := loadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Handle routes
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/users", serveUsersPage)
	http.HandleFunc("/api/users", handleUsers)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server with SSL (HTTPS)
	log.Println("Server started on https://localhost:443")
	err = http.ListenAndServeTLS(":443", certPath, keyPath, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
