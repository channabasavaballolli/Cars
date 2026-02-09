# Initial Performance Baseline (Before Optimization)

**Date**: 2026-02-06
**Status**: Unoptimized Connection Handling
**Load Scenario**: 100 Concurrent Requests via Postman

## CPU Profile Analysis
Based on `pprof` capture of 30 seconds under load.

### 1. Where was the time spent?
*   **GraphQL Engine (~30%)**:
    *   `graphql.executeOperation`, `resolveField`, `Do`.
    *   **Insight**: High CPU cost for parsing and validating GraphQL queries. This is the "cost of doing business" with dynamic queries.
*   **Memory Allocation (~20%)**:
    *   `runtime.mallocgc` (20.45%).
    *   **Insight**: The application generates substantial temporary objects per request, triggering frequent Garbage Collection cycles (`runtime.gcDrain` ~6%).
*   **Database Driver (~14%)**:
    *   `database/sql.(*DB).retry`, `database/sql.withLock`.
    *   **Insight**: Significant overhead from creating/closing connections and fighting for locks. This confirms that **Connection Pooling was missing**.
*   **Idle / Network Wait (~18%)**:
    *   `syscall.Syscall6`.
    *   **Insight**: Time spend waiting for PostgreSQL to reply or for new HTTP requests to arrive.

## Key Bottlenecks Identified
1.  **Connection Churn**: The high presence of `database/sql` internal locking and retries proves that the app was opening a new DB connection for every single request.
2.  **Allocation Rate**: High memory churn is adding latency via GC pauses.

## Next Steps (Planned Optimization)
*   **Enable Connection Pooling**: Set `MaxOpenConns=25` and `MaxIdleConns=25` to eliminate the database driver overhead.
