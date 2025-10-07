package merchant_counter

import (
	"context"
	"sync"
)

type Repository interface {
	GetMerchantCount(ctx context.Context) (int64, error)
}
type Counter struct {
	mu            sync.RWMutex
	merchantCount int64
}

func New(repository Repository) *Counter {
	count, _ := repository.GetMerchantCount(context.Background())
	return &Counter{
		merchantCount: count,
	}
}

func (rc *Counter) Get() int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	count := rc.merchantCount
	return count
}

func (rc *Counter) Increment() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.merchantCount++
}
