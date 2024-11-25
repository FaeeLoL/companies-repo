# Companies Store

## Overview

**Companies Store Microservice** is a RESTful API designed to manage company data. It supports operations like creating, updating, deleting, and retrieving companies. The service is implemented in Go and uses PostgreSQL as the database and Kafka for event-driven messaging.

## Features

- Create, update, delete, and retrieve company information.
- Kafka integration for event publishing on company data changes.
- PostgreSQL as the primary database.
- Docker Compose setup for local development and testing.

Before running the project, ensure the following tools are installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- Optional: [Go](https://go.dev/) (if you want to build and run locally)

---

## Installation

1. Clone the repository:

   ```bash
   git clone <repository_url>
   cd companies-store

2. Create a configuration file:

Copy the example configuration file to customize it for your environment:

    ```bash
    cp configs/config.example.yaml configs/config.yaml

Modify the file as needed (e.g., database credentials, Kafka brokers).

3. Build and run the application using Docker Compose:

    ```bash
    docker-compose up --build

## Usage

### API Endpoints

#### Base URL
http://localhost:8080/api/companies_repo/v1


#### Endpoints

| Method | Endpoint          | Description               |
|--------|-------------------|---------------------------|
| POST   | `/companies`      | Create a new company      |
| PATCH  | `/companies`      | Update an existing company|
| DELETE | `/companies`      | Delete a company          |
| GET    | `/companies`      | Retrieve company details  |

#### Example Request: Create Company

**POST** `/companies`

    ```json
    {
      "id": "01935fed-1a1e-7bb0-8550-109bbcea38a6",
      "name": "Example Co.",
      "description": "A test company",
      "employees_count": 50,
      "registered": true,
      "type": "Corporations"
    }


#### Kafka Events

The service publishes events to Kafka when company data is created, updated, or deleted. Below are the details of the events:

| Event Key         | Description                           | Payload Example                                          |
|--------------------|---------------------------------------|---------------------------------------------------------|
| `create_company`   | Triggered when a company is created   | `{"action":"create_company","data":{"ID":"01935e9a-567c-7cc6-8c38-5b1a3327b43a","Name":"super_company3","Description":"123912089y74t86r2u7iuwfhesdiljk","EmployeesCount":1,"Registered":true,"Type":"NonProfit"},"id_type":"name","identifier":"01935e9a-567c-7cc6-8c38-5b1a3327b43a"}` |
| `update_company`   | Triggered when a company is updated   | `{"action":"update_company","data":{"id":"01935e9a-567c-7cc6-8c38-5b1a3327b43a","employees_count":1230123,"registered":false},"id_type":"uuid","identifier":"01935e9a-567c-7cc6-8c38-5b1a3327b43a"}` |
| `delete_company`   | Triggered when a company is deleted   | `{"action":"delete_company","data":{},"id_type":"uuid","identifier":"01935e9a-567c-7cc6-8c38-5b1a3327b43a"}`               |


## JWT Authentication

The service uses **JWT tokens** for authentication. Below is how the configuration works and how you can test the service with JWT tokens.

### Configuration

JWT-related settings are specified in the configuration file (`configs/config.yaml`):

- **`server.jwt.secret_key`**: Secret key used to sign and verify tokens (HMAC-SHA256).
- **`server.jwt.trusted_issuers`**: A list of trusted issuers. Tokens issued by these will be accepted by the service.

### Testing with JWT.io
You can generate test tokens using JWT.io.

1. Open the JWT.io Debugger.

2. Set the header to:
   ```json 
   {
   "alg": "HS256",
   "typ": "JWT"
   }
   ```

3. Set the payload to something like:

   ```json
   {
   "iss": "your-trusted-issuer",
   "roles": ["admin"],
   "exp": 1716355200
   }
   ```
   iss: Must match one of the trusted_issuers in the configuration.
   roles: Include "admin" for full access.
   exp: Set the expiration timestamp (Unix time).

4. Use the secret key from your configuration file (jwt.secret_key) to sign the token. 

5. Copy the generated token and include it in your API requests as a header:


## Limitations
### Configuration
The service requires a configuration file (configs/config.yaml) to function. This file must be mounted when running the application in Docker.
If you want to run the service using Dockerfile directly (without Docker Compose), you will need to mount the configuration file manually, e.g.:
    
    ```bash
    docker run -v $(pwd)/configs:/app/configs companies-repo-app

### Hardcoded Defaults
The current implementation relies on specific paths and filenames for configuration 
and does not use environment variables for dynamic settings like database credentials or Kafka brokers.

### Testing
Due to lack of time unit and functional tests weren't implemented.

## Direct execution
If you want to run the service directly without Docker, you can build and run the application using Go:

### migrate database 
```bash
go run cmd/main.go migrate-db --config configs/config.yaml
```

### run the application
```bash
go run cmd/main.go http --config configs/config.yaml
```

## To setup dev environment 
### run postgres and kafka
```bash
docker-compose up -f /testing/docker-compose.yml -d
```
### create topic in kafka
```bash
docker exec -it companies-store-kafka kafka-topics --create \
  --bootstrap-server companies-store-kafka:9092 \
  --replication-factor 1 \
  --partitions 1 \
  --topic company_events
```

### check topic created
```bash
docker exec -it companies-store-kafka kafka-topics --list --bootstrap-server companies-store-kafka:9092
```

### run the application
```bash
go run cmd/main.go http --config configs/config-example.yaml
```

## Linter
Execute linters using the following command:
```bash
golangci-lint run -v --timeout 300s
```

## Postman collection
The postman collection might be used for testing the API.