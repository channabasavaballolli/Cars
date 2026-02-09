# Git Workflow Guide

Follow these steps to manage your optimized code using Git and GitHub.

## 1. Initial Setup (If not already done)
If you haven't connected to GitHub yet:
1.  Create a **New Repository** on GitHub (e.g., `car-inventory-service`).
2.  Run the following commands in your terminal:
    ```bash
    git remote add origin https://github.com/YOUR_USERNAME/car-inventory-service.git
    git branch -M main
    ```

## 2. Commit the Optimization Work
You currently have the connection pooling changes and documentation on your local machine. Let's save them.

```bash
# 1. Create a new branch for this feature (Best Practice)
git checkout -b feature/db-connection-pooling

# 2. Add all changes
git add .

# 3. Commit with a descriptive message (include metrics!)
git commit -m "feat(db): implement connection pooling
- Reduces latency from ~300ms to ~1.5ms
- Caps max open connections to 25
- Adds pprof profiling support
"

# 4. Push the branch to GitHub
git push -u origin feature/db-connection-pooling
```

## 3. Merging Changes
1.  Go to your GitHub repository.
2.  You will see a "Compare & pull request" button. Click it.
3.  Review the changes.
4.  Merge the Pull Request into `main`.

## 4. Future Development Workflow
For every new task (e.g., "Add Redis Caching"):
1.  **Sync `main`**: `git checkout main && git pull`
2.  **New Branch**: `git checkout -b feature/add-redis-caching`
3.  **Code & Test**.
4.  **Commit & Push**.
5.  **Pull Request**.
