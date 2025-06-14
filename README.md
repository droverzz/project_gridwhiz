# Project GridWhiz

## Setup

### Prerequisites

- Python 3.x  
- (Specify other dependencies, e.g., Docker, Node.js, Go, k6, if any)

### Installation

```bash
git clone https://github.com/droverzz/project_gridwhiz.git
cd project_gridwhiz
pip install -r requirements.txt
```

---

## Testing

### 1. API Testing with Postman (gRPC)

You can test the API using Postman with the **gRPC protocol**.  
Below are example requests for each endpoint:

<details>
<summary>Expand for endpoint details</summary>

#### 1. Register (No bearer token)
```json
{
  "email": "test@example.com",
  "password": "Mm123456",
  "name": "test"
}
```

#### 2. Login (No bearer token)
```json
{
  "email": "test@example.com",
  "password": "Mm12345"
}
```

#### 3. Logout (Requires bearer token)
```json
{}
```

#### 4. GetUserByID (Requires bearer token)
```json
{
  "id": "684bf307bd1201a8db592ac2"
}
```

#### 5. AddRole (Requires bearer token, only admin)
```json
{
  "target_user_id": "684be197a99e4291f56ab85e",
  "new_role": "user" // or "admin"
}
```

#### 6. ListUsers (Requires bearer token)
```json
{
  "name": "tung2",
  "email": "",
  "page": 1,
  "limit": 3
}
```

#### 7. UpdateProfile (Requires bearer token)
```json
{
  "name": "newname",
  "email": "newemail@example.com"
}
```

#### 8. DeleteProfile (Requires bearer token)
```json
{}
```

#### 9. GeneratePasswordResetToken (Requires bearer token)
```json
{
  "user_id": "684d17c4ef4340af45608ac4"
}
```

#### 10. ResetPassword (Requires bearer token)
```json
{
  "reset_token": "(token from GeneratePasswordResetToken)",
  "new_password": "newStrongPassword13"
}
```
</details>

---

### 2. Load & Stress Testing

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
