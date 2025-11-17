# Memory Bank Quickstart Guide

Welcome to the Claude Code Memory Bank system! This guide will help you get started.

## What is the Memory Bank?

The Memory Bank is a structured documentation system that helps Claude Code maintain context across sessions. Since Claude's memory resets between sessions, these files serve as the project's institutional memory.

## Quick Setup

### 1. Fill Out Core Files

Start by filling out these files with your actual project information:

1. **projectbrief.md** - Start here! This is your foundation.
   - Define what your project is
   - State your goals
   - Document constraints

2. **productContext.md** - Why does your project exist?
   - What problem does it solve?
   - Who are your users?
   - What workflows do they follow?

3. **systemPatterns.md** - How is it built?
   - Document your architecture
   - Note key design patterns
   - Record important technical decisions

4. **techContext.md** - What technologies do you use?
   - List your tech stack
   - Document setup procedures
   - Note dependencies

5. **activeContext.md** - What are you working on now?
   - Current focus
   - Recent changes
   - Next steps

6. **progress.md** - What's done and what's left?
   - Track completed features
   - List in-progress work
   - Document known issues

### 2. Using the System

#### Starting a New Session
Claude will automatically read the memory bank files due to the configuration in `.claude/claude.md`. Just start chatting and Claude will have context.

#### Planning a Feature
Use the `/plan` command:
```
/plan
```
This triggers a structured planning process where Claude will:
1. Review all memory files
2. Ask clarifying questions
3. Draft a comprehensive plan
4. Execute with your approval

#### Updating the Memory Bank
When significant work is done, tell Claude:
```
update memory bank
```
Claude will review ALL memory files and update them with recent changes.

## Best Practices

### Keep It Current
- Update `activeContext.md` frequently (after every significant change)
- Update `progress.md` when features are completed
- Review all files periodically for accuracy

### Be Specific
- Include file paths with line numbers (`src/app.ts:42`)
- Provide code examples in `systemPatterns.md`
- Document the "why" behind decisions

### Maintain Hierarchy
- `projectbrief.md` is the source of truth
- Other files build upon it
- Keep them consistent

## Common Workflows

### Adding a New Feature
1. Use `/plan` to enter planning mode
2. Claude reads memory bank and asks clarifying questions
3. Approve the plan
4. Claude implements with TodoWrite tracking
5. Say "update memory bank" when done

### Understanding the Codebase
1. Start a session (memory loads automatically)
2. Ask Claude questions - it has full context
3. Claude references memory bank files in answers

### Onboarding a Team Member
1. Have them read `projectbrief.md`
2. Then `productContext.md`
3. Then `systemPatterns.md` and `techContext.md`
4. They now understand the project!

## File Structure Reference

```
.claude/
├── claude.md                    # Configuration & hooks
├── commands/
│   └── plan.md                  # Planning mode command
└── memory-bank/
    ├── QUICKSTART.md            # This file
    ├── memory-rules.md          # System documentation
    ├── projectbrief.md          # Foundation document
    ├── productContext.md        # Product vision
    ├── systemPatterns.md        # Architecture
    ├── techContext.md           # Tech stack
    ├── activeContext.md         # Current work
    └── progress.md              # Status tracking
```

## Tips

1. **Start Small**: Fill out basic information first, expand as needed
2. **Iterate**: These files evolve with your project
3. **Be Honest**: Document problems and challenges, not just successes
4. **Use Examples**: Code snippets and examples are very helpful
5. **Link Things**: Reference other memory files when relevant

## Need Help?

- Read `memory-rules.md` for complete documentation
- Check individual file templates for guidance
- Ask Claude to help you fill out any file

## Next Steps

1. Fill out `projectbrief.md` with your project information
2. Work through the other core files
3. Start using `/plan` for your next feature
4. Update memory bank as you work

That's it! You're ready to use the Memory Bank system.
