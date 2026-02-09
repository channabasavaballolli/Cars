# Performance Measurement Guide

## 1. Using Postman (Client-Side)

Postman automatically calculates the Response Time for every request.

### Steps:
1.  Open **Postman**.
2.  Send a Request (e.g., `GET http://localhost:8000/cars`).
3.  Look at the **Response Area** (bottom right of the request pane).
4.  You will see three metrics:
    *   **Status**: (e.g., 200 OK)
    *   **Time**: (e.g., **15ms**) - This is your Response Time.
    *   **Size**: (e.g., 400 B)

### Hover for Details
Hover over the **Time** label to see the breakdown:
*   **Socket Initialization**: Connection setup time.
*   **DNS Lookup**: Domain resolution.
*   **Handshake**: SSL/TLS setup (if HTTPS).
*   **Transfer Start**: Time to first byte (TTFB).
*   **Download**: Time to download the response body.

### Benchmark Testing (Collection Runner)
To get an average over multiple requests:
1.  Save your requests into a **Collection**.
2.  Click the **Runner** button (at the top).
3.  Select your Collection.
4.  Set **Iterations** to `20` (or more).
5.  Click **Run**.
6.  Postman will show the **Average Response Time** for the batch.

---

## 2. Server-Side Logging (Backend)

I have updated the application middleware to log the processing time for every request.

### View Logs:
Run the following command to see the real-time logs:

```powershell
docker-compose logs -f app
```

### Output Example:
```
2026/02/06 10:45:00 Method: GET, URL: /cars, Duration: 2.543ms
2026/02/06 10:45:05 Method: POST, URL: /graphql, Duration: 5.120ms
```

*   **Duration**: This is the time taken *only* by your Go code and Database query (excluding network latency).
*   **Difference**: `Postman Time` - `Server Duration` = `Network Latency` + `Overhead`.

---

## 3. Deep Profiling (pprof)

To see *exactly* which functions are slow (CPU usage) or using too much memory.

### Step A: Start the Profiler
Your app exposes profiling data on port `6060`.
Open: [http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)

### Step B: Generate Traffic (Load Test)
You **must** send requests while profiling, or you will only see idle interactions.
1.  **Postman**: Right-click a Collection -> **Run Collection**.
2.  **Config**: Set **Iterations** to `100`, **Delay** to `0 ms`.
3.  **Don't click Run yet!**

### Step C: Capture the Profile
1.  Run this command in your terminal:
    ```powershell
    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
    ```
2.  **IMMEDIATELY** click **Run** in Postman.
3.  The profiler will record for 30 seconds while Postman hammers the API.

### Step D: Visualize
After the 30 seconds, type:
*   `web`: To see the flow graph (requires Graphviz).
*   `top`: To see the top CPU-consuming functions in text format.

---

## 4. Verifying Optimizations
We have enabled **Connection Pooling** to reduce database overhead.

### Expected Improvements:
1.  **Lower Latency**: Average response time should drop (target < 200ms).
2.  **Less CPU**: `database/sql.retry` and `lock` related functions should disappear from the profile.

### How to Test:
1.  Run the **Postman Collection Runner** again (100 iterations).
2.  Check the **Average Response Time** in the Run Summary.
3.  (Optional) Run a new profile and compare with the old one.
