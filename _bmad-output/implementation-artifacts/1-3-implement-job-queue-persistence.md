# Story : 1.3 — Implement Job Queue Persistence

**Story ID:** 1-3-implement-job-queue-persistence  
**Story Points:** 8  
**Status:** ready-for-dev  
**Created:** 31 janvier 2026  
**Last Updated:** 31 janvier 2026  
**Author:** Epic 1 - Project Foundation & Infrastructure

---

## 📖 Story

As a system, I want to persist job queue state to SQLite, so that jobs survive application restarts.

---

## ✅ Acceptance Criteria

- [ ] **AC1** : Job state transitions (queued → downloading → completed) update the jobs table row
- [ ] **AC2** : `startedAt` timestamp is set when job transitions to "downloading"
- [ ] **AC3** : `completedAt` timestamp is set when job transitions to "completed" or "failed"
- [ ] **AC4** : On application restart, unfinished jobs (status = "queued" or "downloading") are reloaded
- [ ] **AC5** : Reloaded jobs resume from their last checkpoint within 5 seconds
- [ ] **AC6** : Job persistence operations are transactional (all-or-nothing)
- [ ] **AC7** : Concurrent job updates do not cause race conditions or data corruption
- [ ] **AC8** : All tests pass (zero failures, no regressions)

---

## 📋 Tasks

### Task 1: Update Job State Persistence Logic (1.1-1.4)

**Subtasks:**
- **1.1** Update `JobsRepository` to handle state transitions with proper timestamp updates
- **1.2** Ensure `startedAt` is set atomically when state changes to "downloading"
- **1.3** Ensure `completedAt` is set atomically when state changes to "completed" or "failed"
- **1.4** Add transactional updates to prevent partial state changes

### Task 2: Implement Job Recovery on Startup (2.1-2.5)

**Subtasks:**
- **2.1** Create `LoadUnfinishedJobs()` method in JobsRepository
- **2.2** Query jobs with status IN ('queued', 'downloading')
- **2.3** Sort recovered jobs by `createdAt` ASC for FIFO order
- **2.4** Integrate recovery into application startup sequence
- **2.5** Add logging for recovered job count and details

### Task 3: Implement Job Resume Logic (3.1-3.3)

**Subtasks:**
- **3.1** Add checkpoint persistence (track last completed episode)
- **3.2** Implement resume logic to skip already downloaded episodes
- **3.3** Ensure resume happens within 5 seconds of startup

### Task 4: Add Concurrency Safety (4.1-4.3)

**Subtasks:**
- **4.1** Add database-level locking for job updates (SELECT FOR UPDATE or equivalent)
- **4.2** Implement optimistic locking with version field or timestamp comparison
- **4.3** Add tests for concurrent job updates

### Task 5: Write Unit Tests (5.1-5.6)

**Subtasks:**
- **5.1** Test state transition updates (queued → downloading → completed)
- **5.2** Test timestamp setting (startedAt, completedAt)
- **5.3** Test job recovery on startup
- **5.4** Test job resume from checkpoint
- **5.5** Test transaction rollback on failure
- **5.6** Test concurrent updates

### Task 6: Write Integration Tests (6.1-6.3)

**Subtasks:**
- **6.1** Test full lifecycle: create job → restart app → recover job → complete job
- **6.2** Test recovery with multiple unfinished jobs
- **6.3** Test resume performance (must complete within 5 seconds)

---

## 🛠️ Dev Notes

### Architecture Considerations

**Existing Job Repository:**
- Jobs repository already exists at `internal/adapters/sqlite/jobs_repo.go`
- Current schema includes: `id`, `type`, `state`, `progress`, `created_at`, `updated_at`, `params_json`, `result_json`, `error_code`, `error_message`
- Need to verify alignment with story AC schema expectations

**State Transition Safety:**
- Use database transactions for atomic state updates
- Consider SQLite's `IMMEDIATE` transaction mode for write operations
- Implement retry logic for busy database errors

**Job Recovery Strategy:**
1. Load unfinished jobs on startup
2. Re-enqueue into job service
3. Resume from last checkpoint (if supported)
4. Log recovery metrics

**Checkpoint Persistence:**
- Store progress as JSON in `params_json` or `result_json`
- Track: last completed episode number, total episodes, downloaded bytes
- Enable resume by checking progress before downloading each episode

### Testing Strategy

**Unit Tests:**
- Mock database operations
- Test state machine transitions
- Verify timestamp updates
- Test transaction behavior

**Integration Tests:**
- Use real SQLite database (`:memory:` or temp file)
- Test application restart scenario
- Measure recovery time (must be < 5 seconds)

### Performance Considerations

- Recovery query should use index on `state` column (already exists: `idx_jobs_updated_at`)
- Batch loading if > 100 unfinished jobs
- Async recovery to not block application startup

### Error Handling

- Log failed state transitions
- Roll back transactions on error
- Alert if recovery fails (corrupted job data)

---

## 🔗 Dependencies

- Story 1.1 (Initialize SQLite Database Schema) ✅ DONE
- Existing JobsRepository implementation
- Existing Job domain model

---

## 📦 Expected Deliverables

1. Updated `internal/adapters/sqlite/jobs_repo.go` with state transition methods
2. New `LoadUnfinishedJobs()` method in JobsRepository
3. Integration point in `cmd/asd-server/main.go` for job recovery
4. Checkpoint persistence logic in job execution flow
5. Unit tests: `internal/adapters/sqlite/jobs_repo_test.go` (extended)
6. Integration tests: `internal/adapters/sqlite/jobs_integration_test.go` (new file)
7. Updated documentation in story file

---

## Status

**Current Status:** ready-for-dev  
**Progress:** 0/6 major tasks completed (0%)  
**Created:** 31 janvier 2026
