package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomakado/kvevri/store"
)

func TestStore_Set(t *testing.T) {
	t.Run("non-existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
			value = []byte("value")
		)

		store.Set(key, value)

		valueFromStore, ok := store.Get(key)
		assert.True(t, ok)
		assert.Equal(t, value, valueFromStore)
	})

	t.Run("existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
			value1 = []byte("value")
			value2 = []byte("new value")
		)

		store.Set(key, value1)
		store.Set(key, value2)

		valueFromStore, ok := store.Get(key)
		assert.True(t, ok)
		assert.Equal(t, value2, valueFromStore)
	})
}

func TestStore_Get(t *testing.T) {
	t.Run("non-existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
		)

		_, ok := store.Get(key)
		assert.False(t, ok)
	})

	t.Run("existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
			value = []byte("value")
		)

		store.Set(key, value)

		valueFromStore, ok := store.Get(key)
		assert.True(t, ok)
		assert.Equal(t, value, valueFromStore)
	})

	t.Run("expired key", func(t *testing.T) {
		const ttl = 500 * time.Millisecond

		var (
			store = store.New(ttl)
			key   = []byte("key")
			value = []byte("value")
			ctx, cancel = context.WithCancel(context.Background())
		)

		defer cancel()

		go store.StartExpirationWorker(ctx, ttl / 2)

		store.Set(key, value)

		valueFromStore, ok := store.Get(key)
		assert.True(t, ok)
		assert.Equal(t, value, valueFromStore)

		time.Sleep(2 * ttl)

		valueFromStore, ok = store.Get(key)
		assert.False(t, ok)
		assert.Empty(t, valueFromStore)
	})
}

func TestStore_Delete(t *testing.T) {
	t.Run("non-existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
		)

		// should not panic
		store.Delete(key)
	})

	t.Run("existing key", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key   = []byte("key")
			value = []byte("value")
		)

		store.Set(key, value)
		store.Delete(key)

		valueFromStore, ok := store.Get(key)
		assert.False(t, ok)
		assert.Empty(t, valueFromStore)
	})
}

func TestStore_Keys(t *testing.T) {
	t.Run("empty store", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
		)

		keys := store.Keys()
		assert.Empty(t, keys)
	})

	t.Run("non-empty store", func(t *testing.T) {
		var (
			store = store.New(1 * time.Second)
			key1   = []byte("key1")
			value1 = []byte("value1")
			key2   = []byte("key2")
			value2 = []byte("value2")
		)

		store.Set(key1, value1)
		store.Set(key2, value2)

		keys := store.Keys()
		assert.ElementsMatch(t, [][]byte{key1, key2}, keys)
	})
}
