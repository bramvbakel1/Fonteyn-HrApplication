package handlers

import (
	"encoding/json"
	"fmt"
	"hrapplication/internal/azure"
	"hrapplication/internal/utils"
	"net/http"
	"text/template"
)

// Serve the homepage
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/home.html")
	if err != nil {
		http.Error(w, "Error loading homepage", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ServeUsersPage serves the users page.
func ServeUsersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/users.html")
	if err != nil {
		http.Error(w, "Error loading users page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// HandleUsers fetches and responds with a list of users.
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	env, err := utils.LoadEnvironment()
	if err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		return
	}

	accessToken, err := azure.GetAccessToken(env["AZURE_TENANT_ID"], env["AZURE_CLIENT_ID"], env["AZURE_CLIENT_SECRET"])
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	users, err := azure.FetchUsers(accessToken)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// HandleDeleteUser deletes a specific user based on the userID provided in the URL.
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Path[len("/api/users/"):]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	env, err := utils.LoadEnvironment()
	if err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		return
	}

	accessToken, err := azure.GetAccessToken(env["AZURE_TENANT_ID"], env["AZURE_CLIENT_ID"], env["AZURE_CLIENT_SECRET"])
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	err = azure.DeleteUser(accessToken, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
