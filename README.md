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
