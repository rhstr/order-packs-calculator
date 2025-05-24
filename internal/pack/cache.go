package pack

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dgraph-io/ristretto/v2"
	"go.uber.org/zap"
)

// Cache represents a cache for packing results.
type Cache interface {
	// Get retrieves the packing result from the cache for a given order size and pack sizes.
	Get(orderSize int, packs ...int) []Packing
	// Set stores the packing result in the cache for a given order size and pack sizes.
	Set(orderSize int, packs []int, result []Packing)
}

// cache is an implementation of Cache using Ristretto for in-memory caching.
type cache struct {
	ristretto *ristretto.Cache[string, []Packing]
}

// NewInMemoryCache creates a new instance of Cache backed by in-memory storage.
func NewInMemoryCache() (Cache, error) {
	c, err := ristretto.NewCache(&ristretto.Config[string, []Packing]{
		NumCounters: 1e7,       // number of keys to track frequency of (10M).
		MaxCost:     100 << 20, // maximum cost of cache (100M).
		BufferItems: 64,        // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a new Ristretto cache: %w", err)
	}

	return &cache{
		ristretto: c,
	}, nil
}

func (c *cache) Get(orderSize int, packs ...int) []Packing {
	if orderSize <= 0 || len(packs) == 0 {
		return nil
	}

	value, found := c.ristretto.Get(c.key(orderSize, packs))
	if !found {
		zap.L().Info("cache miss", zap.Int("orderSize", orderSize), zap.Ints("packs", packs))

		return nil
	}

	zap.L().Info("cache hit", zap.Int("orderSize", orderSize), zap.Ints("packs", packs))

	return value
}

func (c *cache) Set(orderSize int, packs []int, result []Packing) {
	if orderSize <= 0 || len(packs) == 0 {
		return
	}

	key := c.key(orderSize, packs)

	ok := c.ristretto.Set(key, result, 1)
	if !ok {
		zap.L().Warn("failed to cache the packing", zap.Int("orderSize", orderSize), zap.Ints("packs", packs))
	}
}

// key generates a unique cache key based on the order size and pack sizes.
func (c *cache) key(orderSize int, packs []int) string {
	const delimiter = ":"

	sort.Ints(packs)

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(orderSize))
	sb.WriteString(delimiter)

	for _, pack := range packs {
		sb.WriteString(delimiter)
		sb.WriteString(strconv.Itoa(pack))
	}

	return sb.String()
}
