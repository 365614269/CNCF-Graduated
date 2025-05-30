package cache

import (
	"testing"
)

func TestCacheAddAndGet(t *testing.T) {
	const N = shardSize * 4
	c := New(N)
	c.Add(1, 1)

	if _, found := c.Get(1); !found {
		t.Fatal("Failed to find inserted record")
	}

	for i := range N {
		c.Add(uint64(i), 1)
	}
	for i := range N {
		c.Add(uint64(i), 1)
		if c.Len() != N {
			t.Fatal("A item was unnecessarily evicted from the cache")
		}
	}
}

func TestCacheLen(t *testing.T) {
	c := New(4)

	c.Add(1, 1)
	if l := c.Len(); l != 1 {
		t.Fatalf("Cache size should %d, got %d", 1, l)
	}

	c.Add(1, 1)
	if l := c.Len(); l != 1 {
		t.Fatalf("Cache size should %d, got %d", 1, l)
	}

	c.Add(2, 2)
	if l := c.Len(); l != 2 {
		t.Fatalf("Cache size should %d, got %d", 2, l)
	}
}

func TestCacheSharding(t *testing.T) {
	c := New(shardSize)
	for i := range shardSize * 2 {
		c.Add(uint64(i), 1)
	}
	for i, s := range c.shards {
		if s.Len() == 0 {
			t.Errorf("Failed to populate shard: %d", i)
		}
	}
}

func TestCacheWalk(t *testing.T) {
	c := New(10)
	exp := make([]int, 10*2)
	for i := range 10 * 2 {
		c.Add(uint64(i), 1)
		exp[i] = 1
	}
	got := make([]int, 10*2)
	c.Walk(func(items map[uint64]interface{}, key uint64) bool {
		got[key] = items[key].(int)
		return true
	})
	for i := range exp {
		if exp[i] != got[i] {
			t.Errorf("Expected %d, got %d", exp[i], got[i])
		}
	}
}

func BenchmarkCache(b *testing.B) {
	b.ReportAllocs()

	c := New(4)
	for range b.N {
		c.Add(1, 1)
		c.Get(1)
	}
}
