package domain

import (
	"errors"
	"time"
)

type JobState string

const (
	JobQueued    JobState = "queued"
	JobRunning   JobState = "running"
	JobMuxing    JobState = "muxing"
	JobCompleted JobState = "completed"
	JobFailed    JobState = "failed"
	JobCanceled  JobState = "canceled"
)

func (s JobState) IsTerminal() bool {
	return s == JobCompleted || s == JobFailed || s == JobCanceled
}

type Job struct {
	ID        string
	Type      string
	State     JobState
	Progress  float64
	CreatedAt time.Time
	UpdatedAt time.Time

	// StartedAt is set when job transitions to JobRunning
	StartedAt *time.Time
	// CompletedAt is set when job reaches a terminal state (completed/failed/canceled)
	CompletedAt *time.Time

	ParamsJSON   []byte
	ResultJSON   []byte
	ErrorCode    string
	ErrorMessage string

	// FileListJSON is optional file list metadata stored as JSON
	FileListJSON []byte
}

var ErrInvalidTransition = errors.New("invalid job state transition")

func CanTransition(from, to JobState) bool {
	if from == to {
		return true
	}
	switch from {
	case JobQueued:
		return to == JobRunning || to == JobCanceled || to == JobFailed
	case JobRunning:
		return to == JobMuxing || to == JobCanceled || to == JobFailed
	case JobMuxing:
		return to == JobCompleted || to == JobCanceled || to == JobFailed
	case JobCompleted, JobCanceled, JobFailed:
		return false
	default:
		return false
	}
}
