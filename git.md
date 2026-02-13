# Git Workflow: How to Make Changes

This guide explains the exact steps to follow whenever you want to add a new feature or fix a bug in the future.

## 1. Start Fresh (Sync Main)
Before starting new work, always make sure you are on the `main` branch and have the latest code.

```bash
git checkout main
git pull origin main
```

## 2. Create a Feature Branch
**Never work directly on `main`.** Create a branch for your specific task to keep your work isolated and safe.

```bash
# Example: Adding a users table
git checkout -b feature/add-users-table

# Example: Fixing a bug
git checkout -b fix/login-error
```

## 3. Do Your Work
Write your code, test it, and make your changes. The `main` branch stays safe while you experiment here.

## 4. Save Your Work (Commit)
Stage your changes and save them with a descriptive message.

```bash
# 1. Stage all changes
git add .

# 2. Commit with a message describing WHAT you did
git commit -m "feat: add users table and registration endpoint"
```

## 5. Upload to GitHub (Push)
Send your new branch to GitHub.

```bash
# The '-u' flag links your local branch to the remote one
git push -u origin feature/add-users-table
```

## 6. Merge on GitHub (Pull Request)
1.  Go to your repository on GitHub: [https://github.com/channabasavaballolli/Cars](https://github.com/channabasavaballolli/Cars)
2.  You will see a yellow banner: **"feature/add-users-table had recent pushes"**.
3.  Click **Compare & pull request**.
4.  Review your changes.
5.  Click **Create pull request**.
6.  If everything looks good (and checks pass), click **Merge pull request**.

## 7. Cleanup (Optional)
Once merged, you can delete the branch locally.

```bash
git checkout main
git pull origin main
git branch -d feature/add-users-table
```
