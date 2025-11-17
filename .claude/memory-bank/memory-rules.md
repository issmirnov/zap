# Claude Code Memory Bank System

## Overview

Claude Code's memory resets completely between sessions. This Memory Bank system provides structured documentation that enables effective context preservation and project continuity. The Memory Bank must be read at the start of every task - this is enforced through hooks in claude.md.

## Memory Bank Structure

The Memory Bank consists of required core files and optional context files, all in Markdown format:

```
.claude/
├── claude.md                    # Hooks configuration
└── memory-bank/
    ├── memory-rules.md          # This file - system documentation
    ├── projectbrief.md          # Foundation document (required)
    ├── productContext.md        # Product vision and goals (required)
    ├── activeContext.md         # Current work focus (required)
    ├── systemPatterns.md        # Architecture and patterns (required)
    ├── techContext.md           # Technologies and setup (required)
    ├── progress.md              # Status and achievements (required)
    └── [additional-contexts]/   # Optional feature-specific docs
```

### File Hierarchy

Files build upon each other:
- `projectbrief.md` → Foundation for all other files
- `productContext.md` → Why the project exists (derived from brief)
- `systemPatterns.md` → How it's built (derived from brief)
- `techContext.md` → What technologies (derived from brief)
- `activeContext.md` → Current state (synthesizes product, system, tech)
- `progress.md` → What's done and what's next (tracks active context)

### Core Files (Required)

#### 1. projectbrief.md
- Foundation document that shapes all other files
- Created at project start
- Defines core requirements and goals
- Source of truth for project scope

#### 2. productContext.md
- Why this project exists
- Problems it solves
- How it should work
- User experience goals

#### 3. activeContext.md
- Current work focus
- Recent changes
- Next steps
- Active decisions and considerations

#### 4. systemPatterns.md
- System architecture
- Key technical decisions
- Design patterns in use
- Component relationships

#### 5. techContext.md
- Technologies used
- Development setup
- Technical constraints
- Dependencies

#### 6. progress.md
- What works
- What's left to build
- Current status
- Known issues

### Additional Context
Create additional files/folders within memory-bank/ when they help organize:
- Complex feature documentation
- Integration specifications
- API documentation
- Testing strategies
- Deployment procedures

## Core Workflows

### Starting Any Task
1. **Hooks automatically trigger** - Memory Bank files are read
2. **Verify context** - Ensure understanding is current
3. **Proceed with task** - Use documented patterns and context

### Planning Mode (via /plan command)
1. Read Memory Bank (automatic via hooks)
2. Check if files are complete and current
3. If incomplete: Create plan to establish context
4. If complete: Develop strategy based on existing context
5. Present approach to user

### Execution Mode
1. Check Memory Bank (automatic via hooks)
2. Update documentation as you work
3. Update project intelligence in claude.md if patterns emerge
4. Execute task using documented patterns
5. Document significant changes

## Documentation Updates

Memory Bank updates occur when:
1. Discovering new project patterns
2. After implementing significant changes
3. When user requests **update memory bank** (MUST review ALL files)
4. When context needs clarification
5. After completing major features or milestones

### Update Process
When updating (especially for "update memory bank" command):
1. Review ALL memory bank files
2. Document current state accurately
3. Clarify next steps
4. Update project intelligence in claude.md if needed
5. Focus particularly on activeContext.md and progress.md

## Project Intelligence (claude.md)

The claude.md file serves as a learning journal, capturing:
- Critical implementation paths discovered
- User preferences and workflow patterns
- Project-specific conventions
- Known challenges and solutions
- Evolution of project decisions
- Effective tool usage patterns

### What to Capture in claude.md
Add insights that aren't obvious from code:
- Preferred approaches for this specific project
- Lessons learned from previous attempts
- User's working style and preferences
- Project-specific patterns that work well

The format is flexible - focus on valuable insights that improve future work.

## Automatic Memory Loading

Hooks in claude.md ensure Memory Bank files are automatically read:
- At the start of each new conversation
- When specific commands are triggered
- Before major planning or execution tasks

This automation ensures context is always available without manual intervention.

## Best Practices

1. **Keep It Current**: Update activeContext.md and progress.md frequently
2. **Be Specific**: Vague documentation is worse than no documentation
3. **Show Don't Tell**: Include code examples in systemPatterns.md
4. **Track Decisions**: Document why choices were made, not just what
5. **Maintain Hierarchy**: Each file builds on projectbrief.md
6. **Regular Reviews**: Periodically review all files for accuracy

## Commands

- `/plan` - Enter planning mode with full memory context
- **update memory bank** - Trigger comprehensive review and update of all files

---

**Remember**: After every session reset, the Memory Bank is the only link to previous work. Its accuracy determines effectiveness. Maintain it with precision and clarity.
