# Car Inventory Microservice

A production-ready RESTful and GraphQL Microservice built with Golang, PostgreSQL, and Docker.

## 🚀 Features

*   **Standard Go Layout**: Structured with `models`, `handlers`, `db`, and `utils`.
*   **Database**: PostgreSQL persistence with `lib/pq`.
*   **REST API**: Full CRUD operations for managing cars.
*   **GraphQL API**: Endpoint to query car data flexibly.
*   **Security**:
    *   **Environment Variables**: Secrets management using `godotenv`.
    *   **Authentication**: **Email OTP** & **JWT** for GraphQL. (Legacy API Key for REST).
*   **Validation**: Input validation rules (e.g., Price > 0, Year > 1886).
*   **Observability**: Request logging middleware.
*   **Containerization**: Multi-stage Docker build and Docker Compose orchestration.

## ⚡ Performance Optimization

We achieved a **99% reduction in latency** (from ~300ms to ~1.5ms) by implementing Database Connection Pooling.

| Metric | Before Optimization 🐢 | After Optimization 🚀 | Improvement |
| :--- | :--- | :--- | :--- |
| **Response Time** | 80ms - 500ms+ (Variable) | **~1.5ms** (Consistent) | **50x Faster** |
| **CPU Profile** | Dominated by `syscall.Connect` | Clean Application Logic | **Efficient** |
| **Bottleneck** | DB Connection Handshake | None (at current scale) | **Resolved** |

*See [OPTIMIZATION.md](OPTIMIZATION.md) for full technical details.*

## �️ Tech Stack

*   **Language**: Golang (1.21+)
*   **Router**: Gorilla Mux
*   **Database**: PostgreSQL
*   **GraphQL**: graphql-go
*   **Deployment**: Docker & Docker Compose

## 🏁 Getting Started

### Prerequisites

*   Docker Desktop installed and running.

### 1. Run the Microservice

The entire stack (App + DB) is defined in `docker-compose.yml`.

```powershell
# Build and start services
docker-compose up --build -d
```

*   **App URL**: `http://localhost:8000`
*   **Database**: Port `5432`

### 2. Initialize the Database

Once the containers are running, you need to create the table.

```powershell
# Copy schema to the container
docker cp db/schema.sql cars-db-1:/tmp/schema.sql

# Execute SQL
docker-compose exec db psql -U postgres -d Cars -f /tmp/schema.sql
```

## 🧪 API Documentation & Postman Testing

You can import these details into Postman to test the service.

### 🔐 Authentication & RBAC

The system uses **Role-Based Access Control (RBAC)**:

| Role | Access Level | Login URL |
| :--- | :--- | :--- |
| **User** (Default) | **Read-Only**: Can view the dashboard but cannot modify data. | `/login` |
| **Admin** | **Full Access**: Can Create, Update, and Delete cars. | `/admin-login` |

**Flow:**
1.  **Request Login**: User submits email -> Server sends 6-digit code.
2.  **Verify Login**: User submits email + code -> Server returns **JWT Token**.
3.  **Access Protected Routes**: User sends `Authorization: Bearer <TOKEN>` header.

#### How to Promote a User to Admin
By default, all new users are `user`. To make someone an admin:

1.  Ensure the user has logged in at least once.
2.  Stop the server (if running).
3.  Run the helper script:
    ```powershell
    go run scripts/make_admin.go <user-email>
    ```
4.  Restart the server.

#### 1. Request Login Code
```graphql
mutation {
  requestLogin(email: "your-email@example.com")
}
```
*Result: Check your email for the code.*

#### 2. Verify Code & Get Token
```graphql
mutation {
  verifyLogin(email: "your-email@example.com", code: "123456")
}
```
*Result: Returns a JWT String.*

#### 3. Access Protected Data (Admin Only)
To use mutations like `createCar`, `updateCar`, or `deleteCar`:
*   You must be an **Admin**.
*   Add Header: `Authorization: Bearer <YOUR_JWT_TOKEN>`.

---

### 🔑 Authentication (Legacy REST)
**Headers**:
*   `X-API-Key`: `Infobell` (Only for the REST `DELETE` endpoint)


## GraphQL API Examples

### 1. Create Car (Mutation)
*   **URL**: `http://localhost:8000/graphql`
*   **Method**: `POST`
*   **Body** (GraphQL Mutation):
    ```graphql
    mutation {
        createCar(
            make: "Porsche"
            model: "911"
            year: 2024
            price: 125000.00
            color: "Silver"
            mileage: 0
        ) {
            id
            make
            model
        }
    }
    ```

### 2. Get All Cars (Query)
*   **URL**: `http://localhost:8000/graphql`
*   **Method**: `POST`
*   **Body** (GraphQL Query):
    ```graphql
    query {
        cars {
            id
            make
            model
            year
            price
            color
            mileage
        }
    }
    ```

### 3. Update Car (Mutation)
*   **URL**: `http://localhost:8000/graphql`
*   **Method**: `POST`
*   **Body** (GraphQL Mutation):
    ```graphql
    mutation {
        updateCar(
            id: 1
            price: 130000.00
            mileage: 50
        ) {
            id
            price
            mileage
        }
    }
    ```

### 4. Delete Car (Mutation)
*   **URL**: `http://localhost:8000/graphql`
*   **Method**: `POST`
*   **Body** (GraphQL Mutation):
    ```graphql
    mutation {
        deleteCar(id: 1)
    }
    ```

---

## REST API Examples

### 1. Create a Car (POST)
*   **URL**: `http://localhost:8000/cars`
*   **Method**: `POST`
*   **Body** (JSON):
    ```json
    {
        "make": "Tesla",
        "model": "Model S",
        "year": 2024,
        "price": 89999.99,
        "color": "Red",
        "mileage": 0
    }
    ```
*   **Success Response**: `200 OK` (Returns created Car object)
*   **Validation Error**: Try `year: 1800` to see a `400 Bad Request`.

### 2. Get All Cars (GET)
*   **URL**: `http://localhost:8000/cars`
*   **Method**: `GET`
*   **Response**: List of cars.

### 3. Get Single Car (GET)
*   **URL**: `http://localhost:8000/cars/{id}` (e.g., `/cars/1`)
*   **Method**: `GET`

### 4. Update Car (PUT)
*   **URL**: `http://localhost:8000/cars/{id}`
*   **Method**: `PUT`
*   **Body** (JSON):
    ```json
    {
        "make": "Tesla",
        "model": "Model S Plaid",
        "year": 2024,
        "price": 109999.99,
        "color": "Black",
        "mileage": 500
    }
    ```

### 5. Delete Car (DELETE - Protected)
*   **URL**: `http://localhost:8000/cars/{id}`
*   **Method**: `DELETE`
*   **Headers**: `X-API-Key: Infobell`
*   **Note**: Without the header, you will receive `403 Forbidden`.

## 📂 Project Structure

```
├── db/
│   ├── db.go         # Database connection logic
│   └── schema.sql    # Database schema
├── graph/
│   └── schema.go     # GraphQL schema & resolver
├── handlers/
│   └── cars.go       # REST request handlers
├── models/
│   └── car.go        # Car struct definition
├── utils/
│   └── validation.go # Validation logic
├── .env              # Environment variables
├── docker-compose.yml
├── Dockerfile
└── main.go           # Entry point & Routing
```
