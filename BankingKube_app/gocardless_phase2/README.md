# BankingKube API Development Guide - App Development


## 1. Build API package (portal) on Top of GoCardless/Nordigen API

### Register and Get API Credentials
- Sign up on the Nordigen platform and obtain your API keys or credentials. These are necessary for authenticating your API requests.

### Implement OAuth 2.0 Authentication
- Use the OAuth 2.0 protocol to authenticate and obtain access tokens for making authorized API calls. This process often involves redirecting users to a consent page where they can log in with their bank credentials and grant your application permission to access their banking data.

### Make API Calls
- Use the access token to make authorized API calls to retrieve account information, initiate payments, or access transaction data. This involves sending HTTP requests to the Nordigen API endpoints and handling responses.

### Handle Responses and Errors
- Process the API responses to extract the needed data. Implement error handling to manage any issues that arise during the API interaction.

### Secure Data Handling
- Ensure that all data retrieved from the Nordigen API is handled securely, in compliance with data protection regulations and best practices.








