package domain

import (
	"testing"
	"time"
)

// Test Job state transitions
func TestJob_CanTransition_Valid(t *testing.T) {
	tests := []struct {
		from, to JobState
		valid    bool
	}{
		{JobQueued, JobRunning, true},
		{JobQueued, JobCanceled, true},
		{JobQueued, JobFailed, true},
		{JobRunning, JobMuxing, true},
		{JobRunning, JobCanceled, true},
		{JobRunning, JobFailed, true},
		{JobMuxing, JobCompleted, true},
		{JobMuxing, JobCanceled, true},
		{JobMuxing, JobFailed, true},
		{JobCompleted, JobCompleted, true}, // same state allowed
		{JobCanceled, JobCanceled, true},    // same state allowed
		{JobFailed, JobFailed, true},        // same state allowed
		// Invalid transitions
		{JobQueued, JobCompleted, false},
		{JobQueued, JobMuxing, false},
		{JobRunning, JobQueued, false},
		{JobCompleted, JobRunning, false}, // terminal states can't transition
	}

	for _, tt := range tests {
		result := CanTransition(tt.from, tt.to)
		if result != tt.valid {
			t.Errorf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, result, tt.valid)
		}
	}
}

func TestJobState_IsTerminal(t *testing.T) {
	tests := []struct {
		state    JobState
		terminal bool
	}{
		{JobQueued, false},
		{JobRunning, false},
		{JobMuxing, false},
		{JobCompleted, true},
		{JobFailed, true},
		{JobCanceled, true},
	}

	for _, tt := range tests {
		result := tt.state.IsTerminal()
		if result != tt.terminal {
			t.Errorf("IsTerminal(%q) = %v, want %v", tt.state, result, tt.terminal)
		}
	}
}

// Test Subscription domain model
func TestSubscription_New(t *testing.T) {
	now := time.Now().UTC()
	sub := Subscription{
		ID:                    "sub-1",
		BaseURL:               "https://example.com",
		Label:                 "Test Anime",
		Player:                "auto",
		LastScheduledEpisode:  0,
		LastDownloadedEpisode: 0,
		LastAvailableEpisode:  0,
		NextCheckAt:           now,
		LastCheckedAt:         time.Time{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if sub.ID != "sub-1" {
		t.Errorf("expected ID 'sub-1', got %q", sub.ID)
	}
	if sub.Label != "Test Anime" {
		t.Errorf("expected label 'Test Anime', got %q", sub.Label)
	}
}

// Test Settings domain model
func TestSettings_DefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings.Destination != "videos" {
		t.Errorf("expected destination 'videos', got %q", settings.Destination)
	}
	if settings.OutputNamingMode != OutputNamingLegacy {
		t.Errorf("expected OutputNamingMode 'legacy', got %q", settings.OutputNamingMode)
	}
	if settings.MaxWorkers != 2 {
		t.Errorf("expected MaxWorkers 2, got %d", settings.MaxWorkers)
	}
	if settings.MaxConcurrentDownloads != 4 {
		t.Errorf("expected MaxConcurrentDownloads 4, got %d", settings.MaxConcurrentDownloads)
	}
}

// Test Job domain model
func TestJob_New(t *testing.T) {
	now := time.Now().UTC()
	job := Job{
		ID:           "job-1",
		Type:         "download",
		State:        JobQueued,
		Progress:     0.0,
		CreatedAt:    now,
		UpdatedAt:    now,
		ParamsJSON:   []byte(`{"url": "https://example.com"}`),
		ResultJSON:   []byte{},
		ErrorCode:    "",
		ErrorMessage: "",
	}

	if job.ID != "job-1" {
		t.Errorf("expected ID 'job-1', got %q", job.ID)
	}
	if job.State != JobQueued {
		t.Errorf("expected state 'queued', got %v", job.State)
	}
	if job.Progress != 0.0 {
		t.Errorf("expected progress 0.0, got %v", job.Progress)
	}
}

// Test Job error handling
func TestJob_WithError(t *testing.T) {
	job := Job{
		ID:           "job-1",
		State:        JobFailed,
		ErrorCode:    "TIMEOUT",
		ErrorMessage: "Request timed out",
	}

	if job.ErrorCode != "TIMEOUT" {
		t.Errorf("expected error code 'TIMEOUT', got %q", job.ErrorCode)
	}
	if job.ErrorMessage != "Request timed out" {
		t.Errorf("expected error message 'Request timed out', got %q", job.ErrorMessage)
	}
}

// Test Job progress
func TestJob_Progress(t *testing.T) {
	tests := []struct {
		progress float64
		name     string
	}{
		{0.0, "not started"},
		{0.25, "25% complete"},
		{0.5, "50% complete"},
		{0.75, "75% complete"},
		{1.0, "complete"},
	}

	for _, tt := range tests {
		job := Job{Progress: tt.progress}
		if job.Progress != tt.progress {
			t.Errorf("progress %v: expected %v, got %v", tt.name, tt.progress, job.Progress)
		}
	}
}

// Test Job state constants
func TestJobStateConstants(t *testing.T) {
	if JobQueued != "queued" {
		t.Errorf("JobQueued should be 'queued', got %q", JobQueued)
	}
	if JobRunning != "running" {
		t.Errorf("JobRunning should be 'running', got %q", JobRunning)
	}
	if JobMuxing != "muxing" {
		t.Errorf("JobMuxing should be 'muxing', got %q", JobMuxing)
	}
	if JobCompleted != "completed" {
		t.Errorf("JobCompleted should be 'completed', got %q", JobCompleted)
	}
	if JobFailed != "failed" {
		t.Errorf("JobFailed should be 'failed', got %q", JobFailed)
	}
	if JobCanceled != "canceled" {
		t.Errorf("JobCanceled should be 'canceled', got %q", JobCanceled)
	}
}

// Test Job state transitions are irreversible for terminal states
func TestJobState_TerminalImmutable(t *testing.T) {
	terminalStates := []JobState{JobCompleted, JobFailed, JobCanceled}
	nonTerminalStates := []JobState{JobQueued, JobRunning, JobMuxing}

	for _, terminal := range terminalStates {
		for _, nonTerminal := range nonTerminalStates {
			if CanTransition(terminal, nonTerminal) {
				t.Errorf("terminal state %q should not transition to %q", terminal, nonTerminal)
			}
		}
	}
}
