---
description: Project configuration and memory bank system for Claude Code
alwaysApply: true
---

# Project Configuration

This file contains project-specific intelligence and configuration for Claude Code sessions.

## Memory Bank System

This project uses a Memory Bank system to maintain context across sessions. All memory files are stored in `.claude/memory-bank/` and should be read at the start of every conversation.

### Core Memory Files (Read these in order):
1. **memory-rules.md** - System documentation
2. **projectbrief.md** - Foundation document
3. **productContext.md** - Product vision
4. **systemPatterns.md** - Architecture patterns
5. **techContext.md** - Technical setup
6. **activeContext.md** - Current work focus
7. **progress.md** - Status tracking

## Instructions for Claude

### On Every Session Start:
1. **Read all core memory bank files** in the order listed above
2. Verify understanding of current context
3. Check activeContext.md for immediate priorities
4. Review progress.md to understand what's complete and what's pending

### When User Says "update memory bank":
1. **Review ALL memory bank files** (mandatory - don't skip any)
2. Update files that need changes based on recent work
3. Pay special attention to activeContext.md and progress.md
4. Document any new patterns discovered
5. Update project intelligence in this file if needed

### Planning Mode (via /plan command):
1. Read Memory Bank (automatic)
2. Verify context completeness
3. Ask 4-6 clarifying questions based on findings
4. Draft comprehensive plan
5. Get user approval
6. Implement systematically

## Project Intelligence

This section grows as patterns and preferences are discovered during work on this project.

### Project-Specific Patterns
[Document patterns as they emerge]

### User Preferences
[Document user's working style and preferences]

### Lessons Learned
[Document key insights from work sessions]

### Known Gotchas
[Document tricky areas or common pitfalls]

### Effective Approaches
[Document what works well for this project]

---

## Notes

- Memory Bank files are in `.claude/memory-bank/`
- Templates are provided - fill them out based on actual project
- Update frequently to maintain accuracy
- The Memory Bank is the primary context system for this project

## Architect-Gate: Automated Architecture Review

This project has an automated "Senior Architect in CI" that reviews every PR for architectural consistency and anti-patterns.

**Location**: `.github/claude/` contains prompts and workflow
**Documentation**: See `.github/claude/README.md` for full details

### What It Does
- Runs on every PR automatically
- Reviews changes against Memory Bank and CLAUDE.md
- Checks for 10+ categories of anti-patterns
- Enforces Project-specific architectural invariants
- Posts System Architecture Review (SAR) as PR comment
- Blocks merge if severity ≥ 3 or blockers present

### Key Patterns Enforced
1. **Email Provider Consolidation**: All providers MUST use `EmailProcessingService`
2. **EventBus Contracts**: New events MUST update `EventBus.js` aggregateId derivation
3. **Multi-Tenant Isolation**: All queries MUST filter by `organizationId`
4. **API Routing**: Backend routes have NO `/api` prefix (Caddy handles it)
5. **Deployment Scripts**: Use `deploy-beta.sh` / `deploy-dev.sh`, not manual Docker
6. **Password Hashing**: NEVER manual hashing (User model has hook)
7. **Service Layer**: Follow `Routes → Services → Models` pattern
8. **Async Operations**: Long operations MUST use Redis queues
9. **Docker Rebuild**: Code changes require `--build` flag
10. **Import Pipeline**: Follow 3-phase pattern (Download → Selection → Processing)

### Testing Locally
```bash
# Test prompt composition
./.github/claude/test-prompt-composition.sh

# Review composed prompt
less .architect-gate-test/full-prompt.md

# Test with ACT (requires ANTHROPIC_API_KEY)
act pull_request -s ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY
```

### For PR Authors
- Review SAR comment for specific issues
- Address blocker items before merge
- Push new commits to re-trigger review
- Comment `@architect-run` to manually re-trigger

---
