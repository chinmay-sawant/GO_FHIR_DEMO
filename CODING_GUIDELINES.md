# Go Coding Guidelines & Best Practices 🐹

A set of guidelines and best practices for writing clean, readable, and maintainable Go code for this project.

## Table of Contents
- [Introduction 🚀](#introduction-)
- [Formatting 🎨](#formatting-)
- [Naming Conventions 🏷️](#naming-conventions-️)
- [Comments 💬](#comments-)
- [Error Handling 🚨](#error-handling-)
- [Concurrency ⚡](#concurrency-)
- [Testing 🧪](#testing-)
- [Packages and Project Structure 📦](#packages-and-project-structure-)
- [Dependencies 🔗](#dependencies-)
- [Security 🔒](#security-)
- [Tooling 🛠️](#tooling-)

---

## Introduction 🚀

This document provides coding conventions for the Go programming language. The goal is to improve the readability and maintainability of our codebase. Consistency is key!

---

## Formatting 🎨

- **`gofmt` is non-negotiable.** All Go code in the repository **must** be formatted with `gofmt` (or `goimports`). Most IDEs can be configured to do this automatically on save.
- **Line Length:** Keep lines under 120 characters where possible. This is a soft limit. Readability is more important than strictly adhering to this rule.

---

## Naming Conventions 🏷️

- **Package Names:**
    - Use short, concise, all-lowercase names.
    - Avoid `_` (underscores) or `mixedCaps`.
    - Example: `package http`, `package service`, `package repository`.
- **Variable Names:**
    - Use `camelCase`.
    - Be descriptive but not overly verbose.
    - Short variable names are fine for short-lived variables (e.g., `i` for loop index, `c` for context).
- **Function/Method Names:**
    - Use `PascalCase` for exported functions/methods.
    - Use `camelCase` for unexported (internal) functions/methods.
- **Interfaces:**
    - Prefer single-method interfaces.
    - Name interfaces for what they do. Often, this means adding an `-er` suffix (e.g., `Reader`, `Writer`).
    - Example: `type PatientService interface { ... }` is also acceptable for larger interfaces.
- **Structs:**
    - Use `PascalCase` for exported structs.
    - Use `camelCase` for unexported structs.

---

## Comments 💬

- **Godoc:** All exported types, functions, and constants should have a Godoc comment.
- **Clarity:** Write comments to explain *why* something is done, not *what* is being done. The code should be self-explanatory about the "what".
- **TODOs:** Use `// TODO:` for things that need to be done later. Include a reference to an issue if possible.
    - `// TODO(yourname): Refactor this to use a more efficient algorithm. See #123.`

---

## Error Handling 🚨

- **Always check for errors.** Do not ignore them with `_`.
- **Error messages:**
    - Should not be capitalized.
    - Should not end with punctuation.
    - Provide context using `fmt.Errorf("...: %w", err)`. The `%w` verb is important for error wrapping.
- **Return errors, don't panic.** `panic` should only be used for unrecoverable errors, like a programming mistake that should have been caught during development.

```go
// Good
if err != nil {
    return nil, fmt.Errorf("failed to get patient: %w", err)
}
```

---

## Concurrency ⚡

- **Keep it simple.** Use channels to communicate between goroutines.
- **Avoid sharing memory by communicating.**
- **Use `context.Context`** for cancellation, timeouts, and passing request-scoped values. It should be the first argument to functions that might block.
- **Use `sync.Mutex`** for simple locking when you must share memory.

---

## Testing 🧪

- **Write tests!** Aim for good test coverage for business logic.
- **Table-driven tests** are a great way to test multiple scenarios.
- **Use mocks/stubs** for external dependencies (database, APIs). The `gomock` library is used in this project.
- **Test files** should be named `_test.go` (e.g., `patient_service_test.go`).
- **Test functions** should start with `Test` (e.g., `func TestCreatePatient(t *testing.T)`).

---

## Packages and Project Structure 📦

- We follow a structure inspired by the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).
    - `/cmd`: Main applications for your project.
    - `/internal`: Private application and library code.
    - `/pkg`: Library code that's ok to use by external applications.
    - `/api`: OpenAPI/Swagger specs, JSON schema files.
- **Circular dependencies are a sign of poor package design.** Refactor to avoid them.

---

## Dependencies 🔗

- Use **Go Modules** for dependency management (`go.mod`, `go.sum`).
- Keep dependencies to a minimum. Before adding a new one, consider if the standard library can do the job.
- Run `go mod tidy` to clean up unused dependencies.

---

## Security 🔒

- **Validate all input**, especially from users or external systems.
- **Use parameterized queries** to prevent SQL injection. `gorm` handles this for us, but be mindful.
- **Don't log sensitive information** like passwords or API keys.
- **Manage secrets** using tools like HashiCorp Vault or Consul, not by hardcoding them.

---

## Tooling 🛠️

- **`go vet`:** Run `go vet ./...` to catch suspicious constructs.
- **`golangci-lint`:** A powerful linter that runs many checks. It's highly recommended to integrate it into your workflow.
- **`swag`:** Used for generating Swagger/OpenAPI documentation from comments. Keep the annotations up to date.
- **`mockgen`:** Used for generating mocks for interfaces.

---

Happy Coding! 🎉
