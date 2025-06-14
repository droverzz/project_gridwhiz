# Project GridWhiz

## Setup

### Prerequisites

- Python 3.x  
- (Specify other dependencies, e.g., Docker, Node.js, Go, k6, if any)

### Installation

```bash
git clone https://github.com/droverzz/project_gridwhiz.git
```

---

## API Documentation

- The full gRPC API contract is defined in [proto/auth.proto](proto/auth.proto).
- Example requests and usage for each endpoint are available in [docs/api.md](docs/api.md).

## Architectural Overview

### Tech Stack

- **Go (Golang):**  
  Used for implementing the core business logic and gRPC server. Go offers high performance, simplicity, and strong concurrency, making it ideal for scalable microservices.

- **gRPC (Protocol Buffers):**  
  The service uses gRPC for efficient, strongly-typed, language-agnostic communication between services. All API contracts are defined in Protocol Buffers, which provide clear type safety and automatic code generation for clients in various languages.

- **MongoDB:**  
  Serves as the primary database. MongoDB is a NoSQL database that stores data in flexible, JSON-like documents, suitable for rapid development and horizontal scaling.

- **JWT (JSON Web Token):**  
  Used for stateless authentication. After logging in, users receive a JWT, which must be sent with every authenticated request. Token blacklisting is used to support logout and token invalidation.

---

### MongoDB Schema Overview

#### Collection: users

Example document:
```json
{
  "_id": "ObjectId('684d17c4ef4340af45608ac4')",
  "name": "test",
  "email": "test@example.com",
  "role": "admin",
  "password": "$2a$10$27Ing9Y3yVm9LAoOSwzwROzN/mIMwx50t583Qvygd6iB1M9rpMAS",
  "deleted": false,
  "created_at": "2025-06-14T06:33:40.493+00:00"
}
```
- **_id**: MongoDB ObjectId (Primary key)
- **name**: User name (string)
- **email**: Email address (string, unique)
- **role**: User role ("user" or "admin")
- **password**: Hashed password (string, bcrypt)
- **deleted**: Soft delete status (boolean)
- **created_at**: Account creation timestamp (ISODate string)


---

#### Collection: blacklisted_tokens

Example document:
```json
{
  "_id": "ObjectId('684d1fb11d3c4ff22a36bb')",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3e3NDk5Njc4MzEsInVzZXJfa...",
  "expired_at": "2025-06-15T06:10:31.000+00:00"
}
```
- **_id**: MongoDB ObjectId (Primary key)
- **token**: JWT or refresh token that has been blacklisted (string)
- **expired_at**: Expiry date of the blacklisted token (ISODate string)

---


## Testing

---

###  Load & Stress Testing

The service is designed to handle:
- **~1,000 concurrent users**
- **~100 requests per second**
- **Database capacity for ~100,000 user records**
- **Response times under 200ms for most operations**

To test these requirements:

1. **Seeding test users:**
   ```bash
   cd test
   go run seed_users.go
   ```
   This will populate the database with test users.

2. **Running load tests:**
   ```bash
   # In the test folder
   k6 run loadtest.js
   ```
   This will run a load test to simulate concurrent users and measure system performance.

> **Note:**  
> - Ensure Go and k6 are installed on your system before running the above commands.
> - Review and adjust `seed_users.go` and `loadtest.js` as needed to match your target test scenarios.

---
