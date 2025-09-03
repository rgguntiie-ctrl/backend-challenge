# Backend Golang Coding Test

## Objective
Build a simple RESTful API in Golang that manages a list of users. Use MongoDB for persistence, JWT for authentication, and follow clean code practices.

---

## Requirements

### 1. User Model
Each user should have:
- `ID` (auto-generated)
- `Name` (string)
- `Email` (string, unique)
- `Password` (hashed)
- `CreatedAt` (timestamp)

---

### 2. Authentication

#### Functions
- Register a new user.
- Authenticate user and return a JWT.

#### JWT
- Use JWT for protecting endpoints.
- Use middleware to validate tokens.
- Use HMAC (HS256) with a secret key.

---

### 3. User Functions

- Create a new user.
- Fetch user by ID.
- List all users.
- Update a user's name or email.
- Delete a user.

---

### 4. MongoDB Integration
- Use the official Go MongoDB driver.
- Store and retrieve users from MongoDB.

---

### 5. Middleware
- Logging middleware that logs HTTP method, path, and execution time.

---

### 6. Concurrency Task
- Run a background goroutine every 10 seconds that logs the number of users in the DB.

---

### 7. Testing
Write unit tests

Use Go’s `testing` package. Mock MongoDB where possible.

---

## Bonus (Optional)

- Add Docker + `docker-compose` for API + MongoDB.
- Use Go interfaces to abstract MongoDB operations for testability.
- Add input validation (e.g., required fields, valid email).
- Implement graceful shutdown using `context.Context`.
- **gRPC Version**
  - Create a `.proto` file for `CreateUser` and `GetUser`.
  - Implement a gRPC server.
  - (Optional) Secure gRPC with token metadata.
- **Hexagonal Architecture**
  - Structure the project using hexagonal (ports & adapters) architecture:
    - Separate domain, application, and infrastructure layers.
    - Use interfaces for data access and external dependencies.
    - Keep business logic decoupled from frameworks and DB drivers.

---

## Submission Guidelines

- Submit a GitHub repo or zip file.
- Include a `README.md` with:
  - Project setup and run instructions
  - JWT token usage guide
  - Sample API requests/responses
  - Any assumptions or decisions made

---

## Evaluation Criteria

- Code quality, structure, and readability
- REST API correctness and completeness
- JWT implementation and security
- MongoDB usage and abstraction
- Bonus: gRPC, Docker, validation, shutdown
- Testing coverage and mocking
- Use of idiomatic Go

---------------------------------

## Installation

## Clone the repository:

git clone https://github.com/rgguntiie-ctrl/backend-challenge.git
cd backend-challenge


## Install dependencies:

go mod tidy


## Set up environment variables:

## Create a .env file in the root directory with the following content:

MONGO_URI=mongodb+srv://myuser:mypassword123@cluster0.crdt8fh.mongodb.net/
JWT_SECRET=test-backend-challenge-secret
APP_PORT=3000


## Run the application:

go run ./cmd/backend-api/main.go

## Swagger UI
Open in browser:

http://localhost:3000/docs/index.html