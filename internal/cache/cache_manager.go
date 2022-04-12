package cache

import (
	"context"
)

type HashCacheManager interface {
	Put(ctx context.Context, key string, value map[string]string) error
	Get(ctx context.Context, key string) (map[string]string, error)
}

func NewHashCacheManager(cachedEnabled bool, cacheDir string) (HashCacheManager, error) {
	if cachedEnabled {
		disCache, err := NewDiskCacheManager(cacheDir)
		if err != nil {
			return nil, err
		}
		return disCache, nil
	}
	return &noopCacheManager{}, nil
}

type noopCacheManager struct {
}

func (n noopCacheManager) Put(ctx context.Context, key string, value map[string]string) error {
	return nil
}

func (n noopCacheManager) Get(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}
