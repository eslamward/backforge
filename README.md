# BackForge

BackForge is a Go-based CLI tool that generates CRUD backend applications from a simple YAML schema.

It automatically creates:

* Models
* Handlers
* Repositories
* Routes
* Dependency Injection
* Validation
* Database layer

---

# Table of Contents

* [Installation](#installation)
* [How to use](#how-to-use)
* [How it works](#how-it-works)
* [YAML File Structure Example](#yaml-file-structure-example)
* [Project Generated Structure](#project-generated-structure)
* [Versioning](#versioning)
* [Changelog](#changelog)

---

# Installation

### Option 1: Download Release Binary

Download the latest binary from the releases [page](https://github.com/eslamward/backforge/releases) and place it inside your project folder.

### Option 2: Build from source

```bash
go build -o backforge ./cmd/backforge
```

Then move it to your bin folder.

---

# How to use it

Write an `app.yaml` file with your entities like this: [YAML File Structure Example](#yaml-file-structure-example)

Then run BackForge from the `bin` folder in your OS directory.

```bash
open your terminal in the folder that contains the BackForge tool
backforge build    → generates the code and builds the app
backforge serve    → serves the application on port 8080
```

You can verify the API using:
[http://localhost:8080/api/health](http://localhost:8080/api/health)

This endpoint lists all available APIs for generated entities.

After build, the output will be located in `bin/<os-folder>/` including the `app.db` file.

---

# How it works

```text
app.yaml → parser → generator → filesystem → Go project → build toolchain
```

---

# YAML File Structure Example

## Types you can use:

* integer
* text
* bool
* datetime

## Properties you can use:

* primary
* not_null
* unique
* check (for database only)
* default (for database only)
* auto_increment
* max_value (for numbers)
* min_value (for numbers)
* min_length (for text)
* max_length (for text)

```text
models:
  - name: students
    fields:
      - name: id
        type: integer
        primary: true
        auto_increment: true

      - name: name
        type: text
        not_null: true
        min_length: 4
        check: "length(trim(name)) > 0"

      - name: email
        type: text
        not_null: true
        unique: true
        check: "length(trim(email)) > 0"

      - name: age
        type: integer
        not_null: true
        min_value: 18
        max_value: 70
        check: "age > 5"

  - name: teachers
    fields:
      - name: id
        type: integer
        primary: true
        auto_increment: true

      - name: name
        type: text
        not_null: true
        check: "length(trim(name)) > 0"

      - name: subject
        type: text
        not_null: true

  - name: classes
    fields:
      - name: id
        type: integer
        primary: true
        auto_increment: true

      - name: name
        type: text
        not_null: true

      - name: teacher_id
        type: integer
        not_null: true
        foreign_key:
          model: teachers
          field: id
          on_delete: CASCADE
          on_update: CASCADE

  - name: enrollments
    fields:
      - name: id
        type: integer
        primary: true
        auto_increment: true

      - name: student_id
        type: integer
        not_null: true
        foreign_key:
          model: students
          field: id
          on_delete: CASCADE
          on_update: CASCADE

      - name: class_id
        type: integer
        not_null: true
        foreign_key:
          model: classes
          field: id
          on_delete: CASCADE
          on_update: CASCADE
```

---

# Project Generated Structure

```text
cmd/main.go
internal/
  ├── handler    -- entity handler
  ├── repository -- entity repository
  ├── routes     -- API routing
  ├── app        -- dependency injection container
  ├── models     -- entity models
  ├── database   -- database initialization and table creation
  ├── validate   -- request validation
  └── backerror  -- custom error handling
```

---

# Versioning

BackForge follows semantic versioning:

* Major: breaking changes
* Minor: new features
* Patch: bug fixes

Example:

```
v1.0.0
```

---

# Changelog

## v1.0.0

* Initial release
* YAML parser
* CRUD generator
* SQLite integration
* CLI build and serve commands
