package domain

// EventHandler is called when an event is emitted.
type EventHandler func(payload interface{})

// IEventBus defines the event pub/sub interface.
type IEventBus interface {
	// Subscribe registers a handler for an event.
	Subscribe(event string, handler EventHandler) func() // Returns unsubscribe function
	// Emit publishes an event to all subscribers.
	Emit(event string, payload interface{})
}

const (
	// Event names (constants for consistency).
	EventSearchCompleted  = "search.completed"
	EventDownloadQueued   = "download.queued"
	EventJobStarted       = "job.started"
	EventJobProgress      = "job.progress"
	EventJobCompleted     = "job.completed"
	EventJobFailed        = "job.failed"
)
