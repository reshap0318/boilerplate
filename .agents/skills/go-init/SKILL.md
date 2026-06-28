---
name: go-init
description: Initialize a new Go project from the boilerplate template with mandatory project planning
---

## Role

You are a Go Project Setup Specialist. Your task is to initialize a new Go project using the official boilerplate template, but **only after a complete project plan is in place**.

## When to use me

Use this skill when:

- Starting a new Go project from scratch using the official boilerplate
- The user wants to initialize a project with proper planning (PRD, FSD, Role Matrix, TDD)
- Setting up the initial repository structure, module, and environment
- Migrating or cloning the boilerplate for a new feature or service

## Workflow

### Step 1: Check Project Plan (MANDATORY)

- Verify if the following documents exist in the `docs/` directory:
  - `01_PRD.md`
  - `02_FSD.md`
  - `03_Role_Matrix.md`
  - `04_TDD.md`
- **If ANY are missing:** You MUST invoke `@project-plan` to generate them. Do not proceed to Step 2 until all 4 documents exist.

### Step 2: Confirm Target Path

- Ask the user where to place the project:
  - Root folder (`.`)
  - Subfolder (e.g., `be/`, `backend/`)
- Wait for user confirmation.

### Step 3: Confirm Module Name

- Ask the user for the desired Go module name (e.g., `github.com/user/project`).
- Wait for user confirmation.

### Step 4: Clone Boilerplate

- Clone `https://github.com/reshap0318/boilerplate.git` branch `go` into the target path.
- If cloning into a subfolder, ensure the directory exists first.

### Step 5: Remove Existing Git
- Remove `.git` directory from the cloned project to detach from the boilerplate repository.
- Run `Remove-Item -Recurse -Force -LiteralPath ".git"` on Windows, or `rm -rf .git` on Linux/macOS.

### Step 6: Update Module & Imports

- Edit `go.mod` to use the new module name.
- Recursively find and replace all import paths in `.go` files from `github.com/reshap0318/go-boilerplate` to the new module name.

### Step 7: Setup Environment

- Copy `.env.example` to `.env`.

### Step 8: Verify & Report

- Run `go mod tidy` to ensure dependencies are correct.
- Inform the user that the project is ready.
- Direct them to use skill `@go-development-guide` to start implementing features based on the project plan.
