# Car Service Verification Guide

Follow these steps to start the application and verify all functionality (REST, GraphQL, Security, Validation).

## 1. Start the Environment

Open your terminal in `C:\Users\chann\Desktop\Cars`.

1.  **Clean up old containers (Optional):**
    ```powershell
    docker-compose down -v
    ```

2.  **Start the application:**
    ```powershell
    docker-compose up --build -d
    ```

3.  **Wait for startup:**
    Run `docker-compose logs -f app` and wait until you see:
    > `Successfully connected to database!`
    > `Server starting...`
    (Press `Ctrl+C` to exit logs).

## 2. Initialize the Database

Create the `cars` table in the running PostgreSQL container.

1.  **Copy schema file to container:**
    ```powershell
    docker cp db/schema.sql cars-db-1:/tmp/schema.sql
    ```
    *(Note: If `cars-db-1` is not found, run `docker ps` to check the database container name).*

2.  **Run the SQL:**
    ```powershell
    docker-compose exec db psql -U postgres -d Cars -f /tmp/schema.sql
    ```
    **Output:** `CREATE TABLE`

## 3. Verify REST API & Validation

1.  **Create a Car (Success):**
    ```powershell
    curl.exe -v -X POST http://localhost:8000/cars `
      -H "Content-Type: application/json" `
      -d "{\"make\": \"Tesla\", \"model\": \"Model S\", \"year\": 2024, \"price\": 79999, \"color\": \"Red\", \"mileage\": 0}"
    ```
    **Expect:** `200 OK` and JSON response with ID.

2.  **Validation Test (Fail):**
    Try creating a car with Year 1800 (Too old).
    ```powershell
    curl.exe -v -X POST http://localhost:8000/cars `
      -H "Content-Type: application/json" `
      -d "{\"make\": \"Old\", \"model\": \"Car\", \"year\": 1800, \"price\": 100, \"color\": \"Black\", \"mileage\": 1000}"
    ```
    **Expect:** `400 Bad Request` and error message.

3.  **Get All Cars:**
    ```powershell
    curl.exe http://localhost:8000/cars
    ```

## 4. Verify GraphQL API

1.  **Query Cars via GraphQL:**
    ```powershell
    curl.exe -X POST http://localhost:8000/graphql `
      -H "Content-Type: application/json" `
      -d "{\"query\": \"{ cars { id make model } }\"}"
    ```
    **Expect:** JSON `{"data": {"cars": [...]}}`

## 5. Verify Security

1.  **Delete without Key (Fail):**
    ```powershell
    curl.exe -v -X DELETE http://localhost:8000/cars/1
    ```
    **Expect:** `403 Forbidden`.

2.  **Delete with Key (Success):**
    ```powershell
    curl.exe -v -X DELETE http://localhost:8000/cars/1 `
      -H "X-API-Key: secret-admin-key"
    ```
    **Expect:** `200 OK` and `{"result": "success"}`.

## 6. Verify Middleware Logs

Check that your requests were logged:
```powershell
docker-compose logs app
```
**Expect:** Logs showing `Method: POST, URL: /cars`, `Method: DELETE, URL: /cars/1`, etc.
