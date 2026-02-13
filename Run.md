# How to Run the Car Service Microservice

## Prerequisites
- **Go**: Version 1.22+ installed.
- **PostgreSQL**: Installed and running (or use Docker).
- **Git**: To clone the repository.
- **Environment**: Create a `.env` file (see `.env.example` or instructions below).

## Setup from Scratch

### 1. Clone the Repository
```bash
git clone <repository_url>
cd "Cars - NewAUTH"
```

### 2. Configure Environment
Create a `.env` file in the root directory:
```env
DB_PASSWORD=Channu@4321
API_KEY=Infobell
# Email Configuration (for OTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**Note on Gmail**: You cannot use your regular Gmail password. You must use an **App Password**:
1.  Go to your [Google Account](https://myaccount.google.com/).
2.  Select **Security**.
3.  Under "Signing in to Google", select **2-Step Verification** (Enable it if off).
4.  At the bottom of the 2-Step Verification page, select **App passwords**.
    *(Tip: If you don't see it, search for "App passwords" in the search bar at the top of the Google Account page, or go directly to [https://myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords))*
5.  Select app: **Mail**, device: **Other**, name it "CarApp".
6.  Copy the 16-character password and paste it into `.env` (without spaces).

### 3. Database Setup
You need to create the database and tables.
1.  **Create Database**: Create a database named `Cars` in Postgres.
2.  **Run Migrations**: Execute the SQL in `db/schema.sql` and `db/migrations.sql`.
    ```bash
    psql -U postgres -d Cars -f db/schema.sql
    psql -U postgres -d Cars -f db/migrations.sql
    ```

#### Option 2: Using pgAdmin (GUI)
1.  Open **pgAdmin** and connect to your server.
2.  Right-click "Databases" -> **Create** -> **Database...** -> Name it `Cars`.
3.  Right-click the new `Cars` database -> **Query Tool**.
4.  Open `db/schema.sql` in a text editor, copy the content, paste it into the Query Tool, and hit **Play** (Execute).
5.  Repeat for `db/migrations.sql`.

### 4. Run the Application
You can run it directly with Go or using Docker.

#### Option A: Using Go (Local)
```bash
go mod tidy
go run main.go
```
Server runs at: http://localhost:8000

#### Option B: Using Docker Compose
```bash
docker-compose up --build
```
Server runs at: http://localhost:8000

---

## How to Check Results (Verification)

The API is now protected by **Email OTP Authentication**. You must login to get a token before creating, updating, or deleting cars.

**GraphQL Endpoint**: http://localhost:8000/graphql

### Step 1: Login (Get Token)
1.  **Request Code**:
    ```graphql
    mutation {
      requestLogin(email: "channabasava.bollolli@gmail.com")
    }
    ```
    *Check your email (or server console if SMTP is not configured) for the 6-digit code.*

2.  **Verify Code**:
    ```graphql
    mutation {
      verifyLogin(email: "channabasava.bollolli@gmail.com", code: "123456")
    }
    ```
    *Copy the returned token string (e.g., `eyJ...`).*

### Step 2: Use Token for Protected Operations
Add the token to your HTTP Headers (in Postman or GraphQL Playground):
```json
{
  "Authorization": "Bearer <YOUR_TOKEN>"
}
```

### Step 3: Run Mutations (in Postman)

**URL**: `http://localhost:8000/graphql`
**Method**: `POST`

1.  **Authorization Tab**:
    - Type: **Bearer Token**
    - Token: Paste the token string you got from `verifyLogin`.
    
2.  **Body Tab**:
    - Select **GraphQL**.
    - **Query**:
      ```graphql
      mutation {
        createCar(
          make: "Tesla"
          model: "Model Y"
          year: 2024
          price: 55000.0
          color: "White"
          mileage: 0
        ) {
          id
          make
        }
      }
      ```
    - Click **Send**.

**Update Car**:
```graphql
mutation {
  updateCar(
    id: 1
    price: 52000.0
  ) {
    id
    price
  }
}
```

**Delete Car**:
```graphql
mutation {
  deleteCar(id: 1)
}
```

### Step 4: Public Queries (No Token Needed)
Anyone can list cars without logging in.
```graphql
query {
  cars {
    id
    make
    model
    price
  }
}
```
