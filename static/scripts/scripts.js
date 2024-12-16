// Global variable to store all fetched users
let users = [];

// Function to fetch users from the server and populate the table
async function fetchUsers() {
    try {
        const response = await fetch('/api/users');  // API Endpoint
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        users = await response.json();  // Store fetched users in the users array
        displayUsers(users);  // Display all users in the table
    } catch (error) {
        console.error('Error fetching users:', error);
    }
}

// Function to display the list of users in the table
function displayUsers(usersToDisplay) {
    const tableBody = document.getElementById("users-table-body");
    tableBody.innerHTML = "";

    // Loop through the users array and create a table row for each user
    usersToDisplay.forEach(user => {
        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.displayName}</td>
            <td>${user.userPrincipalName}</td>
            <td>${user.givenName}</td>
            <td>${user.surname}</td>
            <td>${user.roles}</td>
            <td><button onclick="deleteUser('${user.id}')">Delete</button></td>
        `;
        tableBody.appendChild(row);
    });
}

// Function to search users by display name
function searchUsers() {
    const searchInput = document.getElementById("search-input").value.toLowerCase();
    const filteredUsers = users.filter(user => 
        user.displayName.toLowerCase().includes(searchInput)
    );
    displayUsers(filteredUsers);
}

// Function to delete a user
async function deleteUser(userId) {
    if (confirm('Are you sure you want to delete this user?')) {
        try {
            const response = await fetch(`/api/users/${userId}`, {
                method: 'DELETE',
            });

            if (response.ok) {
                alert('User deleted successfully');
                fetchUsers();
            } else {
                alert('Failed to delete the user');
            }
        } catch (error) {
            console.error('Error deleting user:', error);
        }
    }
}

// Function to open the Add User modal
function openModal() {
    const modal = document.getElementById("add-user-modal");
    modal.style.display = "block";  // Show the modal
}

// Function to close the Add User modal
function closeModal() {
    const modal = document.getElementById("add-user-modal");
    modal.style.display = "none";  // Hide the modal
}

// Attach the openModal function to the Add User button
document.getElementById("add-user-button").onclick = openModal;

// Attach the closeModal function to the close button in the modal
document.querySelector(".close-btn").onclick = closeModal;

// Function to handle form submission 
document.getElementById("add-user-form").onsubmit = function(event) {
    event.preventDefault();

    // Collect form data
    const formData = new FormData(event.target);
    const userData = Object.fromEntries(formData.entries());

    // Log the form data to the console
    console.log("New user data:", userData);

    // Close the modal after submission
    closeModal();
};


// Fetch users and display them when the page loads
window.onload = fetchUsers;
