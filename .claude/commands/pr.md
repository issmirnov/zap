---
description: Commit, push, and create or update a pull request
---

# Pull Request Command

Follow these steps to commit changes, push to remote, and create or update a pull request:

## Step 1: Review and Commit Changes

1. Run `git status` and `git diff` to review all changes
2. Analyze the changes and create a meaningful commit message that:
   - Follows the repository's commit message style (check recent commits with `git log`)
   - Clearly describes what was changed and why
   - Is concise but informative
3. Stage all changes with `git add .`
4. Create the commit with the message ending with:

   🤖 Generated with [Claude Code](https://claude.com/claude-code)

   Co-Authored-By: Claude <noreply@anthropic.com>

## Step 2: Push to Remote

1. Check the current branch name with `git branch --show-current`
2. Push the branch to remote with `git push -u origin HEAD` (use `-u` to set upstream if needed)

## Step 3: Check for Existing PR

1. Use `gh pr view --json number,title,body` to check if a PR already exists for this branch
2. If the command returns a PR, proceed to Step 4 (Update PR)
3. If no PR exists (command fails), proceed to Step 5 (Create PR)

## Step 4: Update Existing PR

1. Read the current PR description using the command from Step 3
2. Analyze what new work has been added in the latest commits (use `git log` to see recent commits)
3. Create an updated PR description that:
   - Preserves the original purpose and context
   - Adds a new section or updates existing sections with the latest changes
   - Clearly highlights what's new since the last update
4. Use `gh pr edit --body "$(cat <<'EOF'
[updated description here]
EOF
)"` to update the PR description

## Step 5: Create New PR

1. Analyze all commits on this branch (use `git log main..HEAD` or appropriate base branch)
2. Create a detailed PR with:
   - A clear, descriptive title that summarizes the change
   - A comprehensive body that includes:
     - Summary section: What was changed and why
     - Details section: Key changes made
     - Testing section: How to test or what was tested
     - Any relevant notes or context
3. Use `gh pr create --title "..." --body "$(cat <<'EOF'
## Summary
[summary here]

## Changes
- [change 1]
- [change 2]

## Testing
[testing details]

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"` to create the PR

## Important Notes

- ALWAYS review changes before committing
- NEVER skip hooks or use --no-verify
- Use the main branch as the base unless otherwise specified
- Include the Claude Code attribution in commits and PR descriptions
- Ensure commit messages and PR descriptions are informative and professional
- When updating PRs, clearly indicate what's new to help reviewers
