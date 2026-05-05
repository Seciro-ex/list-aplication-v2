# list-application-v2

A full-stack application for managing users and their purchases, built with Go backend and Vue.js frontend, using PostgreSQL database.

## Features

- User registration and authentication with JWT
- Purchase management (add, edit, delete)
- User profile management
- Swagger API documentation
- Docker support for easy deployment

## Tech Stack

- **Backend**: Go with Gin-like HTTP server, PostgreSQL, JWT authentication
- **Frontend**: Vue.js 3 with Vite, Tailwind CSS
- **Database**: PostgreSQL with migrations
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Go 1.23+
- Node.js 18+
- PostgreSQL (or Docker for containerized setup)

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable

PGADMIN_EMAIL=your_pgadmin_email
PGADMIN_PASSWORD=your_pgadmin_password

JWT_SECRET=your_jwt_secret_key
```

## Installation and Setup

### Option 1: Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd list-application-v2
   ```

2. **Setup PostgreSQL**
   - Install PostgreSQL locally or use Docker:
     ```bash
     docker run --name postgres -e POSTGRES_USER=your_user -e POSTGRES_PASSWORD=your_password -e POSTGRES_DB=your_db -p 5432:5432 -d postgres:15-alpine
     ```

3. **Backend Setup**
   ```bash
   # Install Go dependencies
   go mod download

   # Run database migrations
   go run ./cmd

   # The server will start on http://localhost:8081
   ```

4. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   # Frontend will be available on http://localhost:5173
   ```

### Option 2: Docker Setup

1. **Start services with Docker Compose**
   ```bash
   docker-compose up -d db pgadmin
   # Uncomment app and nginx in docker-compose.yaml if you want to run the full stack in containers
   ```

2. **Run the Go application**
   ```bash
   go run ./cmd
   ```

3. **Run the frontend**
   ```bash
   cd frontend
   npm run dev
   ```

## API Endpoints

All API endpoints are prefixed with `/api`.

### Health Check
**`GET /api/ping`**
- Description: Health check endpoint to verify that the server and database are running
- Response: `Pong! База данных на связи.` (text response)
- Status: `200 OK`

### Authentication

**`POST /api/login`**
- Description: Authenticate user and return JWT token
- Request Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- Response: JWT token
- Status: `200 OK` on success, `401 Unauthorized` on failure

**`POST /api/register`**
- Description: Register a new user in the system
- Request Body:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- Response: `User added successfully!`
- Status: `201 Created` on success, `400 Bad Request` or `500 Internal Server Error` on failure

**`GET /api/me`**
- Description: Get current authenticated user information
- Headers: `Authorization: Bearer <token>`
- Response: User object (without password)
- Status: `200 OK` on success, `401 Unauthorized` if not authenticated

### User Management

**`GET /api/users`**
- Description: Retrieve all registered users from the database (requires authentication)
- Headers: `Authorization: Bearer <token>`
- Response: JSON array of user objects
  ```json
  [
    {
      "id": "integer",
      "username": "string",
      "email": "string",
      "password": "string",
      "created_at": "timestamp"
    }
  ]
  ```
- Status: `200 OK` on success, `500 Internal Server Error` on failure

**`POST /api/edit-user`**
- Description: Edit current user information
- Headers: `Authorization: Bearer <token>`
- Request Body:
  ```json
  {
    "username": "string",
    "email": "string"
  }
  ```
- Response: Success message
- Status: `200 OK` on success, `400 Bad Request` or `500 Internal Server Error` on failure

**`POST /api/delete-user`**
- Description: Delete the current authenticated user
- Headers: `Authorization: Bearer <token>`
- Response: `User deleted successfully!`
- Status: `200 OK` on success, `400 Bad Request` or `500 Internal Server Error` on failure

### Purchase Management

**`POST /api/add-purchase`**
- Description: Add a new purchase/item for the authenticated user
- Headers: `Authorization: Bearer <token>`
- Request Body:
  ```json
  {
    "image_url": "string",
    "cost": "decimal",
    "item_url": "string",
    "item_name": "string"
  }
  ```
- Response: `Purchase added successfully!`
- Status: `201 Created` on success, `400 Bad Request` or `500 Internal Server Error` on failure

**`GET /api/user/purchases`**
- Description: Retrieve all purchases for the authenticated user
- Headers: `Authorization: Bearer <token>`
- Response: JSON array of purchase objects
  ```json
  [
    {
      "id": "integer",
      "image_url": "string",
      "cost": "decimal",
      "item_url": "string",
      "item_name": "string",
      "user_id": "integer"
    }
  ]
  ```
- Status: `200 OK` on success, `500 Internal Server Error` on failure

**`POST /api/edit-purchase`**
- Description: Edit an existing purchase by passing the updated purchase payload in the request body
- Headers: `Authorization: Bearer <token>`
- Request Body:
  ```json
  {
    "id": "integer",
    "image_url": "string",
    "cost": "decimal",
    "item_url": "string",
    "item_name": "string"
  }
  ```
- Response: `Purchase edit successfully!`
- Status: `200 OK` on success, `400 Bad Request` or `500 Internal Server Error` on failure

**`POST /api/delete-purchase`**
- Description: Delete an existing purchase by passing the purchase payload in the request body
- Headers: `Authorization: Bearer <token>`
- Request Body:
  ```json
  {
    "id": "integer"
  }
  ```
- Response: `Purchase deleted successfully!`
- Status: `200 OK` on success, `400 Bad Request` or `500 Internal Server Error` on failure

### API Documentation

**`GET /api/swagger/`**
- Description: Swagger UI for API documentation
- Access: http://localhost:8081/api/swagger/

## Database Migrations

Migrations are automatically run when the server starts. Migration files are located in the `migrations/` directory.

## Building for Production

### Backend
```bash
go build -o bin/main ./cmd
```

### Frontend
```bash
cd frontend
npm run build
```

## Docker Build

To build the Docker image for the Go application:
```bash
docker build -t list-application-v2 .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is licensed under the MIT License.