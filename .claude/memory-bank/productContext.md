# Product Context

> **Derived from**: projectbrief.md
> **Purpose**: Explains why this project exists and what user value it provides

## Problem Statement
### Current Situation
Developers and power users frequently navigate to various web resources throughout their day. Current solutions include:
- **Bookmarks**: Browser-specific, require mouse/menu navigation, become cluttered
- **Search engines**: Require multiple keystrokes, depend on external service
- **Browser history**: Unreliable, requires remembering partial URLs
- **Typing full URLs**: Time-consuming, error-prone, requires memorization

### Pain Points
- Context switching between applications loses browser bookmark access
- Different browsers mean maintaining separate bookmark collections
- Long URLs (e.g., JIRA, GitHub, internal tools) require many keystrokes
- No universal solution that works across all applications
- Chrome shortcuts only work in Chrome omnibox, not in other apps
- SSH, terminal apps, and other tools need full URLs

### Impact
**Time waste**: Typing/navigating to frequently-visited URLs consumes hours weekly. Power users who navigate to dozens of resources daily feel this acutely.

**Context switching cost**: Breaking flow to find bookmarks or copy URLs reduces productivity.

**Cross-platform challenges**: Solutions tied to specific browsers/tools don't work uniformly.

## Solution Overview
### Our Approach
Zap provides a **local HTTP redirect service** that acts as a universal URL shortcut system. Users define shortcuts in a simple YAML config file, then type short codes anywhere that accepts URLs.

**How it works**:
1. Configure shortcuts in `c.yml`: `g: github.com`
2. Type `g/zap` in any application (browser, terminal, Slack, etc.)
3. Zap expands to `https://github.com/issmirnov/zap` and redirects
4. Browser loads the full URL instantly

**Key insight**: By running locally and using HTTP redirects, Zap works universally without browser extensions or OS-specific integrations.

### Key Differentiators
- **Universal**: Works in every application that makes HTTP requests
- **Blazingly fast**: 150k+ requests/second, <15ms latency
- **Hierarchical**: Unlimited nesting (e.g., `g/z/issues/42`)
- **Zero dependencies**: Single binary, no database or external services
- **Hot reload**: Config changes apply instantly without restart
- **Network-wide**: Can serve entire LAN with dnsmasq
- **Powerful**: Supports wildcards, query params, custom schemas, ports

## User Experience
### Target Users
- **Primary**: Software developers navigating between GitHub, docs, tools, localhost services
- **Secondary**: Power users with complex web workflows (researchers, analysts, technical writers)
- **Tertiary**: Technical teams wanting shared shortcuts across organization

### User Workflows

1. **GitHub Navigation**
   - User goal: Quickly access personal repos, issues, PRs
   - Setup: Define `g: github.com` with nested shortcuts for repos
   - Usage: Type `g/zap/issues/42` → redirects to full GitHub issue URL
   - Outcome: Access resources in 10-20 keystrokes instead of 50+

2. **Localhost Development**
   - User goal: Switch between services running on different ports
   - Setup: Define port shortcuts: `api: { port: "3000" }`
   - Usage: Type `api` → redirects to `http://localhost:3000`
   - Outcome: No need to remember port numbers

3. **Search Queries**
   - User goal: Search various sites (Google, StackOverflow, internal wiki)
   - Setup: Use `query` keyword for search expansion
   - Usage: Type `g/golang error handling` → Google search
   - Outcome: Fast access to search results across multiple engines

4. **Team Shared Shortcuts**
   - User goal: Entire team uses consistent shortcuts for internal tools
   - Setup: Shared `c.yml` in git repo, deployed via Ansible
   - Usage: All team members type `jira/PROJ-123` for tickets
   - Outcome: Consistent navigation patterns, reduced onboarding friction

### User Interface Principles
- **Invisible by design**: Users never see Zap's interface - only results
- **Fail fast**: Invalid shortcuts return errors immediately
- **Forgiving**: Hot reload allows quick config fixes
- **Transparent**: `/varz` endpoint shows current configuration
- **Low friction**: YAML config file is simple, human-readable

## Product Requirements
### Must Have (Completed)
- HTTP redirect service with configurable shortcuts
- YAML configuration with hierarchical structure
- Hot reload without service restart
- Cross-platform binaries (Linux, macOS)
- Performance capable of handling developer workload (100k+ QPS)
- Comprehensive documentation and examples

### Should Have (Completed)
- Query parameter expansion for search workflows
- Port number shortcuts for localhost
- Wildcard path preservation
- /etc/hosts integration
- Installation via package managers (Homebrew)
- Automated deployment (Ansible role, systemd)
- Health check and config dump endpoints

### Could Have (Future Enhancements)
- Metrics endpoint for usage statistics
- Config validation web UI
- Import from browser bookmarks
- Shared config repository/registry
- Windows native support with MSI installer
- Mobile companion app

## Success Metrics
- **Performance**: Achieves 150k+ QPS ✓
- **Adoption**: Homebrew formula available ✓
- **Reliability**: Test coverage ~50% of codebase ✓
- **Maintainability**: Active dependency updates ✓
- **Documentation**: Comprehensive README and CONTRIBUTING ✓
- **Ecosystem**: Ansible role available in Galaxy ✓

Current indicators show strong product-market fit within target audience of technical power users.

## Evolution
### Project Maturity
The "zap2" folder name is a workspace convention for parallel work, not a version indicator. The current codebase shows significant maturity:
- Professional CI/CD setup demonstrates established development practices
- Comprehensive testing indicates thorough quality standards
- Multiple deployment methods show community-driven ecosystem growth
- Performance optimization shows real-world usage has informed design

### Recent Evolution
- **Go 1.24 upgrade**: Staying current with language improvements
- **Dependency updates**: Regular Dependabot PRs show active maintenance
- **Architecture stability**: Core design hasn't changed, indicating strong foundation
- **Community tools**: Ansible role, Homebrew formula show ecosystem growth
