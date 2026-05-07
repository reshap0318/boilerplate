---
name: ai-coding-workflow
description: Structured AI coding workflow with plan approval, GitHub issues, branching, PR creation, and issue closure
---

## What I do
- Create a `plan.md` file before any coding begins
- Wait for explicit user approval at EVERY phase before proceeding
- Create GitHub issues via MCP (if available)
- Create feature branches from `master`
- Implement code according to the approved plan
- Commit and push changes to GitHub
- Open Pull Requests to `master`
- Close GitHub issues after PR approval

## When to use me
Use this skill whenever the user asks you to write code, implement features, fix bugs, or refactor. **NEVER skip this workflow** — always start with `plan.md`.

---

## CRITICAL: STOP & WAIT Rule

**After EVERY phase, you MUST STOP and WAIT for the user to respond.**
**NEVER auto-proceed to the next phase without explicit user confirmation.**

The user will inform you when they approve or want to continue. Do not assume approval. Do not chain phases together.

---

## Workflow Steps

### Phase 1: Planning

1. **Analyze the request** — Understand what the user wants
2. **Create `plan.md`** in the project root with this structure:

```markdown
# Plan: {Feature/Bug Name}

## Objective
{Brief description of what will be done}

## Scope
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

## Implementation Details

### Step 1: {Step name}
- Files to create/modify:
  - `path/to/file.go` — description
- Key changes: ...

### Step 2: {Step name}
- Files to create/modify:
  - `path/to/file.go` — description
- Key changes: ...

## Branch Name
`feature/{short-description}`

## Estimated Impact
- Files created: N
- Files modified: N
- Risk level: low/medium/high
```

3. **Present the plan to the user** and say:
   > "Plan sudah dibuat. Apakah kamu approve? Ketik 'approve' atau 'ya' untuk melanjutkan, atau berikan feedback jika ada yang perlu diubah."

4. **STOP and WAIT** — do NOT proceed to Phase 2 until the user explicitly says they approve

5. If user requests changes, **update `plan.md`** and ask again

---

### Phase 2: GitHub Issue (if MCP GitHub available)

**ONLY start this phase after user has approved the plan in Phase 1.**

1. **Check if GitHub MCP is available** by testing the `gh` tool or MCP GitHub tools
2. **If available**, create a GitHub issue:
   - Title: `[Plan] {Feature/Bug Name}`
   - Body: Content from `plan.md`
   - Labels: `enhancement` or `bug` (based on context)
3. **Store the issue number** — you will need it later to close the issue
4. **If not available**, skip this phase and notify the user
5. **After done, STOP and WAIT** — say:
   > "Issue GitHub sudah dibuat (#{number}). Ketik 'lanjut' untuk membuat branch baru."

---

### Phase 3: Branch Creation

**ONLY start this phase after user says 'lanjut' or confirms to continue.**

1. **Ensure working tree is clean** — stash or commit any uncommitted changes first
2. **Create and checkout new branch:**

```bash
git checkout master
git pull origin master
git checkout -b feature/{short-description}
```

3. **Push branch to remote:**

```bash
git push -u origin feature/{short-description}
```

4. **STOP and WAIT** — say:
   > "Branch `feature/{short-description}` sudah dibuat dan di-push. Ketik 'lanjut' untuk mulai coding."

---

### Phase 4: Implementation

**ONLY start this phase after user says 'lanjut' or confirms to continue.**

1. **Follow `plan.md` step by step** — implement each task in order
2. **After each step**, verify the code compiles:

```bash
go build ./...
go vet ./...
```

3. **Update `plan.md`** — check off completed tasks:

```markdown
- [x] Task 1 — completed
- [ ] Task 2
```

4. **Commit incrementally** — commit after each logical step:

```bash
git add .
git commit -m "feat: {step description}"
```

5. **Push changes:**

```bash
git push origin feature/{short-description}
```

6. **When ALL tasks are done, STOP and WAIT** — say:
   > "Semua task di plan sudah selesai. Apakah ada perubahan yang perlu ditambahkan? Jika tidak, ketik 'lanjut' untuk membuat PR."

---

### Phase 5: Pull Request

**ONLY start this phase after user confirms no more changes and says 'lanjut'.**

1. **Create a PR:**

```bash
gh pr create \
  --base master \
  --head feature/{short-description} \
  --title "{Feature/Bug Name}" \
  --body-file plan.md
```

2. **If `gh` is not available**, inform the user to create the PR manually and provide the branch name
3. **Notify the user** with the PR URL and say:
   > "PR sudah dibuat: {PR_URL}. Silakan review dan approve PR secara manual di GitHub. Setelah PR di-merge, ketik 'merged' untuk menutup issue dan cleanup."
4. **STOP and WAIT** — do NOT proceed to Phase 6 until user says 'merged'

---

### Phase 6: Issue Closure

**ONLY start this phase after user confirms PR is merged by saying 'merged'.**

1. **Close the GitHub issue:**

```bash
gh issue close #{issue-number}
```

2. **If MCP GitHub is not available**, remind the user to close the issue manually
3. **Clean up** — delete the local and remote feature branch:

```bash
git checkout master
git branch -d feature/{short-description}
git push origin --delete feature/{short-description}
```

4. **Delete `plan.md`** or archive it
5. Say:
   > "Workflow selesai. Issue #{number} sudah ditutup, branch sudah dihapus, dan plan.md sudah diarsipkan."

---

## Rules

- **NEVER skip Phase 1** — always create `plan.md` first
- **NEVER auto-proceed** — STOP and WAIT after every phase until user confirms
- **NEVER code without approval** — wait for explicit user approval
- **ALWAYS commit with conventional commits** (`feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`)
- **ALWAYS verify code compiles** before committing
- **NEVER push directly to `master`** — always use feature branches
- **ALWAYS update `plan.md`** progress as you work
- If any step fails, **STOP and notify the user** — do not proceed blindly

## Error Handling

- If `git` commands fail, show the error and ask the user for help
- If `gh` commands fail, provide manual instructions for the user
- If build/vet fails, fix the issue before committing
- If PR creation fails, provide the user with the exact command to run manually
