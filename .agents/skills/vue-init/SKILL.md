---
name: vue-init
description: Initialize a new Vue project from the boilerplate template with mandatory project planning
---

## Role
You are a Vue Project Setup Specialist. Your task is to initialize a new Vue project using the official boilerplate template, but **only after a complete project plan is in place**.

## When to use me

Use this skill when:
- Starting a new Vue project from scratch using the official boilerplate
- The user wants to initialize a project with proper planning (PRD, FSD, Role Matrix, TDD)
- Setting up the initial repository structure and environment
- Migrating or cloning the boilerplate for a new frontend project

## Workflow

### Step 1: Check Project Plan (MANDATORY)
- Verify if the following documents exist in the `docs/` directory:
  - `01_PRD.md`
  - `02_FSD.md`
  - `03_Role_Matrix.md`
  - `04_TDD.md`
- **If ANY are missing:** Invoke the `project-plan` skill to generate them. Do not proceed to Step 2 until all 4 documents exist.

### Step 2: Confirm Target Path
- Ask the user where to place the project:
  - Root folder (`.`)
  - Subfolder (e.g., `fe/`, `frontend/`, `web/`)
- Wait for user confirmation.

### Step 3: Clone Boilerplate
- Clone the boilerplate into the target path using the `vue` branch:
  ```bash
  git clone --branch vue https://github.com/reshap0318/boilerplate.git <target-path>
  ```
- If cloning into a subfolder, ensure the directory exists first.

### Step 4: Remove Existing Git
- Remove `.git` directory from the cloned project to detach from the boilerplate repository.
- Run `Remove-Item -Recurse -Force -LiteralPath ".git"` on Windows, or `rm -rf .git` on Linux/macOS.

### Step 5: Update Project Config
Update the following files to match the new project:
- **`index.html`** — change the `<title>` tag to the project name
- **`package.json`** — change the `"name"` field to the project name (lowercase, kebab-case)

### Step 6: Setup Environment
- Copy `.env.example` to `.env`:
  - **Linux/macOS:** `cp .env.example .env`
  - **Windows:** `Copy-Item .env.example .env`
- Open `.env` and fill in the required values:
  - `VITE_API_BASE_URL` — base URL of the backend API (e.g., `http://localhost:8080/api`)
  - `VITE_APP_NAME` — display name of the application

### Step 7: Install Dependencies & Verify
- Run `yarn install` to install all dependencies.
- Run `yarn build` to verify the project builds successfully.
- Inform the user that the project is ready and provide the next steps based on the project plan.
