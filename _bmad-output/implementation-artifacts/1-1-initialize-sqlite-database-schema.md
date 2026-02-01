# Story : 1.1 — Initialize SQLite Database Schema

**Story ID:** 1-1-initialize-sqlite-database-schema  
**Story Points:** 5  
**Status:** ready-for-dev  
**Created:** 31 janvier 2026  
**Last Updated:** 31 janvier 2026  
**Author:** Epic 1 - Project Foundation & Infrastructure

---

## 📖 Story

As a developer, I want to initialize the SQLite database with core tables, so that the application can persist settings and job queue data.

---

## ✅ Acceptance Criteria

- [x] **AC1** : SQLite database file created at `data/asd.db` on first application startup
- [x] **AC2** : `settings` table created with schema: `id`, `key`, `value`, `updated_at`
- [x] **AC3** : `jobs` table created with schema: `id`, `animeId`, `status`, `episodes`, `createdAt`, `startedAt`, `completedAt`, `errorMessage`
- [x] **AC4** : PRAGMA integrity_check returns "ok" status after table creation
- [x] **AC5** : Default settings inserted (download_path, quality, hls_enabled) with sensible defaults
- [x] **AC6** : Database initialization is idempotent (running twice produces same result, no errors)
- [x] **AC7** : All tables include proper primary keys and timestamps with current UTC time
- [x] **AC8** : Zero test failures and all existing tests pass (no regressions)

---

## 🎯 Tasks/Subtasks

### Task 1 : Create Database Initialization Module
- [ ] **1.1** : Create `internal/adapters/sqlite/schema.go` file
- [ ] **1.2** : Define `InitDatabase(dbPath string) error` function
- [ ] **1.3** : Implement idempotent CREATE TABLE IF NOT EXISTS logic
- [ ] **1.4** : Set up proper SQLite schema with primary keys and constraints
- [ ] **1.5** : Test database initialization with unit tests

### Task 2 : Create Settings Table Schema
- [ ] **2.1** : Define settings table with `id`, `key`, `value`, `updated_at` columns
- [ ] **2.2** : Implement default settings insertion (download_path, quality, hls_enabled)
- [ ] **2.3** : Add validation for settings table structure
- [ ] **2.4** : Test settings table operations

### Task 3 : Create Jobs Table Schema
- [ ] **3.1** : Define jobs table with all required columns: `id`, `animeId`, `status`, `episodes`, `createdAt`, `startedAt`, `completedAt`, `errorMessage`
- [ ] **3.2** : Set up proper indexes on `id` and `status` for query performance
- [ ] **3.3** : Implement constraints (NOT NULL for required fields)
- [ ] **3.4** : Test jobs table operations

### Task 4 : Integrate with Application Startup
- [ ] **4.1** : Call `InitDatabase()` during application bootstrap in `cmd/asd-server/main.go`
- [ ] **4.2** : Handle database initialization errors with proper logging
- [ ] **4.3** : Verify idempotency (run twice, verify no errors)
- [ ] **4.4** : Test application startup with fresh database

### Task 5 : Database Integrity Validation
- [ ] **5.1** : Implement PRAGMA integrity_check validation after table creation
- [ ] **5.2** : Return error if integrity check fails
- [ ] **5.3** : Log successful initialization with details
- [ ] **5.4** : Add tests for integrity validation

### Task 6 : Unit Tests & Validation
- [ ] **6.1** : Create `internal/adapters/sqlite/schema_test.go`
- [ ] **6.2** : Test InitDatabase() with fresh database
- [ ] **6.3** : Test idempotency (InitDatabase() called twice)
- [ ] **6.4** : Test default settings are properly inserted
- [ ] **6.5** : Test all tables exist with correct schema
- [ ] **6.6** : Test PRAGMA integrity_check passes
- [ ] **6.7** : Run all existing tests to ensure no regressions

---

## 📝 Dev Notes

### Architecture & Patterns

**Database Setup**:
- SQLite file stored at `data/asd.db` (relative to working directory)
- Use parameterized queries to prevent SQL injection
- Apply migrations via idempotent CREATE TABLE IF NOT EXISTS
- Timestamps in UTC, stored as ISO 8601 strings or Unix timestamps

**Schema Design**:
- `settings`: Key-value store for application configuration
  - `id INTEGER PRIMARY KEY AUTOINCREMENT`
  - `key TEXT UNIQUE NOT NULL`
  - `value TEXT NOT NULL`
  - `updated_at DATETIME DEFAULT CURRENT_TIMESTAMP`

- `jobs`: Queue of download jobs
  - `id TEXT PRIMARY KEY` (UUID or similar)
  - `animeId TEXT NOT NULL`
  - `status TEXT NOT NULL` (values: queued, downloading, completed, failed, cancelled)
  - `episodes TEXT NOT NULL` (e.g., "1-12" or JSON)
  - `createdAt DATETIME NOT NULL`
  - `startedAt DATETIME`
  - `completedAt DATETIME`
  - `errorMessage TEXT`

**Default Settings**:
- `download_path`: `/home/user/Downloads/anime` (or platform-specific default)
- `quality`: `720p`
- `hls_enabled`: `true`

### Dependencies & Libraries

- `database/sql` (Go stdlib)
- `github.com/mattn/go-sqlite3` (SQLite driver)
- Standard library logging

### Known Issues & Constraints

- SQLite file path must be accessible and writable
- First-time initialization may take a few hundred milliseconds
- Database must be closed properly before application exit

### Testing Strategy

- Unit tests using in-memory `:memory:` database
- Test table creation, idempotency, schema validation
- Test default settings insertion
- Verify PRAGMA integrity_check passes
- No external service dependencies needed

---

## 📝 Dev Agent Record

**Implementation Status:** Not started  
**Tests Created:** 0  
**Files Modified:** 0  
**Decisions Made:** None yet  

---

## 📦 File List

*To be updated as implementation progresses*

---

## 📋 Change Log

### Session 1 Changes (31 janvier 2026)

*No changes yet*

---

## Status

**Current Status:** ready-for-dev  
**Progress:** 0/6 major tasks completed (0%)  
**Created:** 31 janvier 2026
