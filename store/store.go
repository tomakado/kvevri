package store

import (
	"context"
	"log"
	"time"

	cache "github.com/go-pkgz/expirable-cache/v2"
)

type Store struct {
	cache cache.Cache[string, []byte]
	ttl   time.Duration
}

func New(ttl time.Duration) *Store {
	return &Store{
		cache: cache.NewCache[string, []byte]().WithTTL(ttl),
		ttl:   ttl,
	}
}

func (s *Store) Set(key, value []byte) {
	s.cache.Set(string(key), value, s.ttl)
}

func (s *Store) Get(key []byte) ([]byte, bool) {
	return s.cache.Get(string(key))
}

func (s *Store) Delete(key []byte) {
	s.cache.Invalidate(string(key))
}

func (s *Store) Keys() [][]byte {
	var (
		keys     = s.cache.Keys()
		byteKeys = make([][]byte, 0, len(keys))
	)

	for _, key := range keys {
		byteKeys = append(byteKeys, []byte(key))
	}

	return byteKeys
}

func (s *Store) StartExpirationWorker(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	s.cache.DeleteExpired()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("deleting expired keys")
			s.cache.DeleteExpired()
		}
	}
}
