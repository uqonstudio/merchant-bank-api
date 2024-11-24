# Merchant-Bank API

This is a simple RESTful API for handling login, payment, and logout functionalities between a merchant and a bank. The application uses JSON files to simulate customer, merchant, and history data.

## Prerequisites

- Ensure that the `customer.json`, `history.json`, `merchant.json` and `payment.json` files exist in the `database` directory.
- The `JWT_LIFE_TIME`, `JWT_ISSUER_NAME`, `JWT_KEY` environment variable must be set for authentication purposes.
- The application uses the Gin framework and requires Go modules for dependency management.

## API Endpoints

### 1. Login

- **Endpoint**: `/api/auth/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }

- **Response**:
- ***200 OK***: Login successful
- ***401 Unauthorized***: Invalid credentials

### 2. Payment

- **Endpoint**: /api/payment-merchant
- **Method**: POST
- **Auth**: Bearer Token
- **Request Body**:
  ```json
  {
    "customer_id": "string",
    "merchant_id": "string",
    "amount": "float",
    "transaction_id": "string"
  }

- **Response**:
- ***200 OK***: Login successful
- ***401 Unauthorized***: Invalid credentials

### 3. Logout

- **Endpoint**: /api/auth/logout
- **Method**: POST
- **Request Body**:
  ```json
  {
    "customer_id": "string"
  }

- **Response**:
    - **200 OK**: Logout successful
    - **401 Unauthorized**: Unauthorized or invalid customer

### 4. Create Customer

- **Endpoint**: api/customers
- **Method**: POST
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }

- **Response**:
- **201 Created:**: The customer was created successfully.
- **400 Bad Request**: The request payload is malformed or missing
- **405 Method Not Allowed**: The request method is not POST.
- **500 Internal Server Error**: An error occurred on the server while processing the request.

### 5. Customer

- **Endpoint**: /api/customers/
- **Method**: GET

## Setup Instructions

### Prerequisites
- Go 1.16 or later installed on your machine.
- Git installed for cloning the repository.

### Installation

#### 1. Clone the Repository:
 ```json
    git clone https://github.com/uqonstudio/merchant-bank-api.git/
    cd merchant-bank-api
```
#### 2. Initialize Go Modules:
```json
    go mod tidy
```

#### 3. Run the Application:
```json
    go run main.go
```

#### 4. Access the API:
- The server will start on http://localhost:8080.
- Use a tool like Postman or curl to interact with the API endpoints.

## Testing
- Implement unit tests using Go's testing package.
- Run tests with:
    ```
    go test ./...
    ```

## Deployment
- Deploy the application on a server or cloud platform that supports Go applications.
- Ensure the JSON files are accessible and writable for logging history.

## Security Considerations
- Use HTTPS for secure communication.
- Implement proper authentication mechanisms (e.g., JWT) for production.
- Validate and sanitize all inputs to prevent injection attacks.
