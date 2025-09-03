
## Installation

## Clone the repository:

- git clone https://github.com/rgguntiie-ctrl/backend-challenge.git
- cd backend-challenge


## Install dependencies:

go mod tidy


## Set up environment variables:

## Create a .env file in the root directory with the following content:

- MONGO_URI=mongodb+srv://myuser:mypassword123@cluster0.crdt8fh.mongodb.net/
- JWT_SECRET=test-backend-challenge-secret
- APP_PORT=3000


## Run the application:

go run ./cmd/backend-api/main.go

## Swagger UI
Open in browser:

http://localhost:3000/docs/index.html

- You can click the Authorize button in Swagger and enter Bearer <JWT token> to authenticate.

- Protected endpoints require a valid JWT token in the Authorization header.