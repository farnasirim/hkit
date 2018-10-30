package hkit

import (
	"context"
)

// CacheParams is the metadata passed to the cache engine.
// Instantiate *only* using the New... Functions and not &CacheParams{}
type CacheParams struct {
	cacheKey string
}

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
