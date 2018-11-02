package hkit

import (
	"context"
)

// InMemoryCacheService implements a trivial goroutine safe cache service using
// memory allocated dynamically in *this process*
type InMemoryCacheService struct {
}

// TryGet implements TryGet of CacheService
func (c *InMemoryCacheService) TryGet(ctx context.Context, cacheParams *CacheParams) (interface{}, error) {
	return nil, nil
}

// Set implements Set of CacheService
func (c *InMemoryCacheService) Set(ctx context.Context, cacheParams *CacheParams, value interface{}) error {
	return nil
}
