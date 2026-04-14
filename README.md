# list-aplication-v2

A Go application for managing users and their purchases with PostgreSQL database integration.

## ENV FILE STRUCTURE
```
DB_USER= ***
DB_PASSWORD= ***
DB_NAME= ***
DATABASE_URL= postgres://user:password@host:port/dbname

PGADMIN_EMAIL= ***
PGADMIN_PASSWORD= ***
```

## API Endpoints

### Health Check
**`GET /ping`**
- Description: Health check endpoint to verify that the server and database are running
- Response: `Pong! База данных на связи.` (text response)
- Status: `200 OK`

### User Management

**`POST /register`**
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

**`GET /users`**
- Description: Retrieve all registered users from the database
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

### Purchase Management

**`POST /add-purchase`**
- Description: Add a new purchase/item for a user
- Request Body:
  ```json
  {
    "image_url": "string",
    "cost": "decimal",
    "item_url": "string",
    "item_name": "string",
    "user_id": "integer"
  }
  ```
- Response: `Purchase added successfully!`
- Status: `201 Created` on success, `400 Bad Request` or `500 Internal Server Error` on failure

**`GET /purchases`**
- Description: Retrieve all purchases from the database
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