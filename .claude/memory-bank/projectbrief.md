# Project Brief

> **Foundation Document**: This file shapes all other memory bank files. Update this first when project scope or goals change.

## Project Name
Zap - High-Performance URL Shortcut Expander

## Purpose
Zap is a fast, local HTTP redirect service that transforms short codes into full URLs. Instead of typing long URLs or managing bookmarks, users define universal web shortcuts in a YAML config file. Typing shortcuts like `g/z` expands to `github.com/issmirnov/zap`. Works in any application with network access.

**GitHub**: github.com/issmirnov/zap

## Core Goals
1. **Maximize user efficiency** by reducing keystrokes needed to access frequently-visited web resources
2. **Deliver exceptional performance** - handle 150k+ requests per second with minimal latency
3. **Maintain simplicity** - single binary, single config file, zero runtime dependencies
4. **Provide flexibility** - support arbitrary URL structures, nesting, and custom workflows
5. **Ensure reliability** - hot reload without downtime, graceful error handling

## Scope
### In Scope
- HTTP redirect service that expands shortcuts to full URLs
- Hierarchical shortcut definitions with unlimited nesting depth
- Hot reload of configuration without service restart
- Query parameter expansion for search workflows
- Port shortcuts for localhost development
- Custom URL schema support (chrome://, etc.)
- Wildcard matching for arbitrary path preservation
- /etc/hosts integration for DNS-level shortcuts
- Cross-platform support (macOS, Linux)
- Multiple deployment methods (Homebrew, systemd, Ansible)

### Out of Scope
- User authentication or multi-tenant support (local-first design)
- Web UI for configuration management (YAML file is interface)
- Cloud/SaaS offering (self-hosted only)
- URL analytics or tracking
- Database or persistent storage (config file is state)
- Windows native support (though binary may work)

## Success Criteria
- Achieves 150k+ QPS sustained performance
- Config changes apply within seconds via hot reload
- Single binary deployment with no runtime dependencies
- Installation available via major package managers (Homebrew)
- Comprehensive test coverage maintained (~50% of codebase)
- Active maintenance with timely dependency updates
- Clear documentation for installation and configuration

## Constraints
- **Technical**:
  - Go 1.24+ required for development
  - Must run on localhost or private network (not internet-facing)
  - Config file must be valid YAML
  - Requires port 80 access for direct use, or reverse proxy setup

- **Time**:
  - No specific deadlines (maintained open source project)
  - Hot reload required to avoid deployment downtime

- **Resources**:
  - Solo maintainer (issmirnov)
  - Open source contributors for features/fixes
  - CI/CD via GitHub Actions (free tier)

- **External**:
  - Browser behavior with non-standard schemas (Chrome blocks custom protocols)
  - Operating system /etc/hosts file access requires elevated permissions
  - Package manager approval processes (Homebrew, Ansible Galaxy)

## Key Stakeholders
- **Primary Users**: Developers and power users who frequently navigate to various web resources
- **Maintainer**: issmirnov (original author)
- **Contributors**: Open source community via GitHub
- **Ecosystem**: Users of Ansible role, Homebrew formula

## Context
The project folder is named "zap2" as a parallel workspace convention, not because this is version 2. The project represents a mature codebase with established patterns and proven performance. The project demonstrates professional software engineering practices including comprehensive testing, CI/CD automation, semantic versioning, and cross-platform release management. The codebase is compact (~1,260 LOC) but highly effective, with roughly equal parts production code and tests.
