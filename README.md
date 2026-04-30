# Git Contribution Guidelines


## 1. Branch Naming Convention

We use a prefix-based naming system. Branches should be named in **kebab-case** (lowercase with hyphens).

**Format:**  
`<type>/<short-description>`

### Allowed Types:
- **feat**: New features or functional additions.
- **fix**: All bug fixes and patches.
- **refactor**: Code changes that neither fix a bug nor add a feature.
- **docs**: Documentation only changes.
- **style**: Changes that do not affect the meaning of the code.
- **test**: Adding missing tests or correcting existing tests.
- **chore**: Changes to the build process or auxiliary tools.

**Examples:**
- `feat/user-registration-api`
- `fix/login-button-overlap`
- `refactor/optimize-database-queries`

---

## 2. Semantic Commit Messages

Format: `<type>(<scope>): <subject>`  
*`<scope>` is optional.*

### Type Definitions:
- **feat**: A new feature for the user.
- **fix**: A bug fix for the user.
- **docs**: Changes to the documentation.
- **style**: Formatting, missing semi-colons, etc; no production code change.
- **refactor**: Refactoring production code (e.g., renaming a variable).
- **test**: Adding missing tests, refactoring tests; no production code change.
- **chore**: Updating build tasks, package manager configs, etc.

### Examples:
- **feat(auth):** implement JWT token refresh
- **fix(ui):** resolve navbar mobile responsiveness
- **docs:** update installation instructions in readme
- **style:** fix indentation in logger.go
- **test:** add unit tests for payment gateway

---
