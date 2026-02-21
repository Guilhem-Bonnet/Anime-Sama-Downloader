package ports

import (
	"context"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// ISchedulerService defines the interface for scheduling subscription checks.
// This abstraction allows dependency injection and enables mocking for unit tests.
type ISchedulerService interface {
	// ScheduleNextCheck schedules the next check time for a subscription.
	// Updates when the subscription should be checked next based on the interval.
	//
	// Parameters:
	//   ctx: context for cancellation and timeouts
	//   subscriptionID: unique ID of the subscription
	//   interval: how often to check (e.g., 1*time.Hour, 24*time.Hour)
	//
	// Returns:
	//   error: if scheduling fails
	ScheduleNextCheck(ctx context.Context, subscriptionID string, interval time.Duration) error

	// GetNextCheckTime retrieves when a subscription should be checked next.
	//
	// Parameters:
	//   ctx: context for cancellation and timeouts
	//   subscriptionID: unique ID of the subscription
	//
	// Returns:
	//   time.Time: the scheduled next check time
	//   error: if retrieval fails or subscription not found
	GetNextCheckTime(ctx context.Context, subscriptionID string) (time.Time, error)

	// GetDueSubscriptions retrieves subscriptions that are due for checking.
	// Returns all subscriptions where NextCheckAt <= now, up to a limit.
	//
	// Parameters:
	//   ctx: context for cancellation and timeouts
	//   now: current time for comparison
	//   limit: maximum number of subscriptions to return
	//
	// Returns:
	//   []domain.Subscription: subscriptions that are due for checking
	//   error: if retrieval fails
	GetDueSubscriptions(ctx context.Context, now time.Time, limit int) ([]domain.Subscription, error)
}
