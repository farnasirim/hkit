package hkit

import (
	"context"
	"net/http"
)

// CacheParams is the metadata passed to the cache engine.
// Instantiate *only* using the New... Functions and not &CacheParams{}
type CacheParams struct {
	cacheKey string
}

// NewCacheParamsByKey returns cache
func NewCacheParamsByKey(cacheKey string) *CacheParams {
	return &CacheParams{
		cacheKey: cacheKey,
	}
}

// CacheParamsFunc is the signature of the function that generates metadata to
// be used for caching the output of an `http.HandlerFunc`. You'll use it to
// generate a cache key to for example cache /resource and /resource/ the same
// way. You may set cache timeouts, eviction policies, etc.
type CacheParamsFunc func(context.Context, *http.Request) *CacheParams

// CacheService defines the api that should be provided by any cache backend
// service. An example would be `RedisCacheService`, implementing caching
// strategies based on the passed `*CacheParams`, using redis as storage backend
type CacheService interface {
	// TryGet will ask the CacheService for the cached value represented by
	// cacheParams, which will also determine what the cacheBackend should do
	// with it after a possible hit (update the lru, refresh timeout, etc.)
	TryGet(ctx context.Context, cacheParams *CacheParams) (interface{}, error)

	// Set will ask the cache backend to associate the value to the supplied
	// *cacheParams struct
	Set(ctx context.Context, cacheParams *CacheParams, value interface{}) error
}
