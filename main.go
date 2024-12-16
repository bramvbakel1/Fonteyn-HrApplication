package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// User struct represents a user in Microsoft Graph API
type User struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	UserPrincipalName string `json:"userPrincipalName"`
	FirstName         string `json:"givenName"` // First Name
	LastName          string `json:"surname"`   // Last Name
	Role              string `json:"role"`
}

// Load environment variables and return them
func loadEnv() (string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", "", "", fmt.Errorf("error loading .env file")
	}

	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	graphAPIURL := os.Getenv("GRAPH_API_URL")

	return tenantID, clientID, clientSecret, graphAPIURL, nil
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
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading homepage", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Serve the users page
func serveUsersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/users.html")
	if err != nil {
		http.Error(w, "Error loading users page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Handle fetching users from Microsoft Graph API
func handleUsers(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	tenantID, clientID, clientSecret, _, err := loadEnv()
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Value []User `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Error decoding response body: %v", err)
	}

	return result.Value, nil
}

// Handle adding a new user
func addUser(w http.ResponseWriter, r *http.Request) {
	// Decode the user data from the request body
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Error decoding user data", http.StatusBadRequest)
		return
	}

	// Load environment variables
	tenantID, clientID, clientSecret, _, err := loadEnv()
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

	// Make API call to create the new user
	err = createUserInGraphAPI(newUser, accessToken)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
}

// Function to create a new user in Microsoft Graph API
func createUserInGraphAPI(newUser User, accessToken string) error {
	url := "https://graph.microsoft.com/v1.0/users"
	userData := map[string]interface{}{
		"displayName":       newUser.DisplayName,
		"userPrincipalName": newUser.UserPrincipalName,
		"givenName":         newUser.FirstName,
		"surname":           newUser.LastName,
		"jobTitle":          newUser.Role,
	}

	// Marshal the user data into JSON
	body, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("error marshaling user data: %v", err)
	}

	// Create request to add user
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to add user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add user: %s", resp.Status)
	}

	return nil
}

// Handle deleting a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	userID := strings.TrimPrefix(r.URL.Path, "/api/users/")

	// Ensure the user ID is not empty
	if userID == "" {
		http.Error(w, "User ID is missing", http.StatusBadRequest)
		return
	}

	// Load environment variables
	tenantID, clientID, clientSecret, _, err := loadEnv()
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

	// Make API call to delete the user
	err = deleteUserFromGraphAPI(userID, accessToken)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusNoContent)
}

// Function to delete a user via Microsoft Graph API
func deleteUserFromGraphAPI(userID, accessToken string) error {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete user: %s", resp.Status)
	}

	return nil
}

// Main server initialization
func main() {
	// Handle routes
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/users", serveUsersPage)
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/users/", deleteUser)
	http.HandleFunc("/api/users/add", addUser)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
