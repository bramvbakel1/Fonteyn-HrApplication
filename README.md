# Go Application with Azure Integration

This Go application integrates with Azure services (Key Vault and Microsoft Graph API) to manage users.

Create a .env file in the root directory of your project with the following contents:
AZURE_TENANT_ID=<your-tenant-id>
AZURE_CLIENT_ID=<your-client-id>
AZURE_CLIENT_SECRET=<your-client-secret>
AZURE_KEYVAULT_URL=<your-keyvault-url>
AZURE_CERT_NAME=<your-cert-name>

go mod tidy

docker build -t myapp .

docker run -d -p 443:443 --name myapp-container myapp
