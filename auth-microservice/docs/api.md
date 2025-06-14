# API Testing with Postman (gRPC)

You can test the API using Postman with the **gRPC protocol**.  
Below are example requests for each endpoint:

---

### 1. Register (No bearer token)
```json
{
  "email": "test@example.com",
  "password": "Mm123456",
  "name": "test"
}
```

### 2. Login (No bearer token)
```json
{
  "email": "test@example.com",
  "password": "Mm12345"
}
```

### 3. Logout (Requires bearer token)
```json
{}
```

### 4. GetUserByID (Requires bearer token)
```json
{
  "id": "684bf307bd1201a8db592ac2"
}
```

### 5. AddRole (Requires bearer token, only admin)
```json
{
  "target_user_id": "684be197a99e4291f56ab85e",
  "new_role": "user" // or "admin"
}
```

### 6. ListUsers (Requires bearer token)
```json
{
  "name": "tung2",
  "email": "",
  "page": 1,
  "limit": 3
}
```

### 7. UpdateProfile (Requires bearer token)
```json
{
  "name": "newname",
  "email": "newemail@example.com"
}
```

### 8. DeleteProfile (Requires bearer token)
```json
{}
```

### 9. GeneratePasswordResetToken (Requires bearer token)
```json
{
  "user_id": "684d17c4ef4340af45608ac4"
}
```

### 10. ResetPassword (Requires bearer token)
```json
{
  "reset_token": "(token from GeneratePasswordResetToken)",
  "new_password": "newStrongPassword13"
}
```