# Optimization Strategies & Layers

This document details the performance optimization techniques applied to the Car Inventory Service, categorized by architectural layer.

## 1. Database Layer (PostgreSQL)
**Technique:** Connection Pooling

*   **Problem:** The application was opening a new TCP connection and performing a handshake (SSL/Auth) for every single HTTP request.
*   **Metric (Before):** Latency ~80ms-500ms. CPU dominated by `database/sql.Open` and `syscall.Connect`.
*   **Optimization:** Configured `database/sql` to maintain a pool of "warm" connections.
    ```go
    // db/db.go
    DB.SetMaxOpenConns(25)    // Cap max concurrency
    DB.SetMaxIdleConns(25)    // Keep them open for reuse
    DB.SetConnMaxLifetime(5 * time.Minute) // Recycle slowly
    ```
*   **Metric (After):** Latency ~1.5ms. CPU overhead for DB connection is near-zero.
*   **Trade-offs:** Consumes memory for idle connections. Requires tuning based on available RAM and Postgres `max_connections` limit.

## 2. Service Layer (Go API)
**Technique:** Concurrency Limiting (via Pool Size)

*   **Problem:** Unbounded concurrency can crash downstream services (Postgres) under load spikes.
*   **Optimization:** The `MaxOpenConns(25)` setting acts as a semaphore. If 100 requests arrive, 25 proceed, and 75 wait in a fast in-memory queue.
*   **Benefit:** Prevents "thundering herd" issues. Ensures consistent response times even under heavy load.

## 3. Observability Layer (Profiling)
**Technique:** On-Demand Profiling (`pprof`)

*   **Strategy:** We do not guess bottlenecks. We use `net/http/pprof` to capture production traffic.
*   **Workflow:**
    1.  Expose `localhost:6060/debug/pprof`.
    2.  Run Load Test (Postman).
    3.  Capture 30s Profile: `go tool pprof ...`
    4.  Analyze "Top" consumers.
*   **Impact:** Identified that 60% of CPU was wasted on DB Handshakes, preventing premature optimization of GraphQL resolvers.

## 4. Future Optimizations (Backlog)
*   **Caching:** Implement Redis for `GetCars` read queries to bypass DB entirely.
*   **Query Optimization:** Indexing `make` and `year` columns if search filters are added.
*   **GraphQL Complexity Limiting:** Prevent deeply nested queries from exhausting server resources.
