---
description: Update memory bank and compress old data if needed
---

# Memory Bank Update Command

Follow this comprehensive process to update and maintain the memory bank:

## Step 1: Review All Memory Bank Files

Read and analyze ALL memory bank files in order:
1. `.claude/memory-bank/memory-rules.md`
2. `.claude/memory-bank/projectbrief.md`
3. `.claude/memory-bank/productContext.md`
4. `.claude/memory-bank/systemPatterns.md`
5. `.claude/memory-bank/techContext.md`
6. `.claude/memory-bank/activeContext.md`
7. `.claude/memory-bank/progress.md`

As you review, identify:
- Outdated information that needs updating
- Completed items that should be archived
- New patterns or insights to document
- Information that's grown too long and needs compression

## Step 2: Check for Compression Opportunities

Look for data that should be compressed or archived:
- **activeContext.md**: Move completed work to progress.md, keep only active items
- **progress.md**: Summarize long lists of completed items, group by milestones
- **systemPatterns.md**: Consolidate similar patterns, remove obsolete approaches
- **techContext.md**: Archive deprecated dependencies, consolidate setup instructions

## Step 3: Update Files Systematically

Update each file that needs changes:

### Priority Updates (do these first):
1. **activeContext.md** - Update current focus, move completed work to progress.md
2. **progress.md** - Add completed items, compress old accomplishments, update status

### Secondary Updates (if needed):
3. **systemPatterns.md** - Document new patterns discovered, remove outdated approaches
4. **techContext.md** - Update dependencies, setup steps, infrastructure changes
5. **productContext.md** - Update if product direction or requirements changed
6. **projectbrief.md** - Only update if core project scope changed

## Step 4: Compress Old Data

When files become too long or cluttered:

### For activeContext.md:
- Keep only current and next 1-2 sprints of work
- Move completed items to progress.md
- Summarize old decisions, keep only key points

### For progress.md:
- Group old completed items by major milestones
- Summarize lists that exceed 15-20 items
- Keep recent work detailed, compress older work

### For systemPatterns.md:
- Consolidate similar patterns into comprehensive examples
- Archive deprecated patterns in a separate section
- Keep examples concise but complete

### For techContext.md:
- Remove deprecated dependencies
- Consolidate repetitive setup instructions
- Archive old infrastructure details if no longer relevant

## Step 5: Update Project Intelligence

Review `.claude/claude.md` and update the Project Intelligence section if you discovered:
- New project-specific patterns that work well
- User preferences or workflow patterns
- Effective approaches for this project
- Known gotchas or challenges
- Lessons learned from recent work

## Step 6: Summarize Changes

After updating, provide a clear summary:
- Which files were updated and why
- What was compressed or archived
- Key new information added
- Current state of the project based on updated memory bank

## Compression Guidelines

**Compress when**:
- activeContext.md exceeds 200 lines
- progress.md has more than 20 completed items in recent work
- systemPatterns.md has redundant or very similar patterns
- Any file becomes difficult to scan quickly

**Don't compress**:
- Critical technical details needed for implementation
- Important architectural decisions and their rationale
- Active work in progress
- Recent changes (last 2-4 weeks)

**Compression techniques**:
- Summarize lists: "Implemented 15 API endpoints for user management" instead of listing all 15
- Group by theme: Combine related completed items
- Archive to separate files: Move detailed specs to feature-specific docs
- Remove redundancy: Consolidate duplicate information

## Important Notes

- ALWAYS review ALL files, even if you think they haven't changed
- Focus on accuracy - don't guess or assume
- Be specific with updates - include file references and line numbers
- Document WHY decisions were made, not just WHAT was done
- Keep the memory bank actionable and useful, not just historical
- After compression, verify that no critical information was lost

---

**Goal**: Maintain a clean, current, and actionable memory bank that provides maximum context with minimum noise.
