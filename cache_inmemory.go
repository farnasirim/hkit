package hkit

// InMemoryCacheService implements a trivial goroutine safe cache service using
// memory allocated dynamically in *this process*
type InMemoryCacheService struct {
}

func (c *InMemoryCacheService) TryGet(ctx context.Context, cacheParams *CacheParams) (interface{}, error) {

}

func (c *InMemoryCacheService) Set(ctx context.Context, cacheParams *CacheParams, value interface{}) error {

}
