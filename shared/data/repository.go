package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"soliant-mock-api/shared/redis"
)

// Entity is the contract any model must satisfy to be stored by Repository.
// Implementations should define methods on the pointer receiver.
type Entity interface {
	GetID() string
	SetID(id string)
}

// Repository is a generic CRUD store backed by Redis.
//
// T is the value type of the model (e.g. Shift), and PT is *T constrained to
// implement Entity (so we can mutate the ID from inside the repository).
type Repository[T any, PT interface {
	*T
	Entity
}] struct {
	cache    redis.CacheRepository
	prefix   string
	notFound error
	seed     func() []T
}

// Config holds the parameters needed to build a Repository.
type Config[T any] struct {
	Cache    redis.CacheRepository
	Prefix   string // e.g. "shifts:"
	NotFound error
	Seed     func() []T
}

// New returns a new generic Repository.
func New[T any, PT interface {
	*T
	Entity
}](cfg Config[T]) *Repository[T, PT] {
	return &Repository[T, PT]{
		cache:    cfg.Cache,
		prefix:   cfg.Prefix,
		notFound: cfg.NotFound,
		seed:     cfg.Seed,
	}
}

// NotFound returns the configured not-found error for this repository.
func (r *Repository[T, PT]) NotFound() error { return r.notFound }

// SeedIfEmpty stores the default items only when no items exist yet.
func (r *Repository[T, PT]) SeedIfEmpty() error {
	if r.seed == nil {
		return nil
	}

	keys, err := r.cache.GetAllKeys(r.prefix)
	if err != nil {
		return fmt.Errorf("read existing %s: %w", r.prefix, err)
	}

	if len(keys) > 0 {
		return nil
	}

	for _, item := range r.seed() {
		item := item
		p := PT(&item)
		if err := r.save(p); err != nil {
			return fmt.Errorf("seed %s%s: %w", r.prefix, p.GetID(), err)
		}
	}

	return nil
}

// GetAll returns every item stored in Redis under this repository's prefix.
func (r *Repository[T, PT]) GetAll() ([]T, error) {
	keys, err := r.cache.GetAllKeys(r.prefix)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(keys))
	for _, key := range keys {
		id := r.extractID(key)
		if id == "" {
			continue
		}
		item, err := r.GetByID(id)
		if err != nil {
			if errors.Is(err, r.notFound) {
				continue
			}
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

// GetByID returns a single item by its id.
func (r *Repository[T, PT]) GetByID(id string) (*T, error) {
	raw, err := r.cache.Get(r.key(id))
	if err != nil {
		return nil, r.notFound
	}

	if len(raw) == 0 {
		return nil, r.notFound
	}

	var item T
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, fmt.Errorf("unmarshal %s%s: %w", r.prefix, id, err)
	}

	return &item, nil
}

// Create persists a new item. If the item's ID is empty a new one is generated.
func (r *Repository[T, PT]) Create(item T) (*T, error) {
	p := PT(&item)

	if p.GetID() == "" {
		nextID, err := r.nextID()
		if err != nil {
			return nil, err
		}
		p.SetID(nextID)
	} else {
		exists, err := r.cache.Exists(r.key(p.GetID()))
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("%s%s already exists", r.prefix, p.GetID())
		}
	}

	if err := r.save(p); err != nil {
		return nil, err
	}

	return &item, nil
}

// Update applies a partial (deep-merge) update to the item identified by id.
// Only fields present in the patch overwrite the stored item; nested objects
// are merged key by key.
func (r *Repository[T, PT]) Update(id string, patch map[string]any) (*T, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(existing)
	if err != nil {
		return nil, err
	}

	var current map[string]any
	if err := json.Unmarshal(raw, &current); err != nil {
		return nil, err
	}

	delete(patch, "id")
	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return nil, err
	}

	var updated T
	if err := json.Unmarshal(merged, &updated); err != nil {
		return nil, err
	}
	PT(&updated).SetID(PT(existing).GetID())

	if err := r.save(PT(&updated)); err != nil {
		return nil, err
	}

	return &updated, nil
}

// Delete removes an item by id.
func (r *Repository[T, PT]) Delete(id string) error {
	exists, err := r.cache.Exists(r.key(id))
	if err != nil {
		return err
	}
	if !exists {
		return r.notFound
	}

	return r.cache.Delete(r.key(id))
}

func (r *Repository[T, PT]) save(p PT) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return r.cache.Set(r.key(p.GetID()), payload, 0)
}

func (r *Repository[T, PT]) key(id string) string { return r.prefix + id }

func (r *Repository[T, PT]) nextID() (string, error) {
	keys, err := r.cache.GetAllKeys(r.prefix)
	if err != nil {
		return "", err
	}

	maxNum := 0
	for _, key := range keys {
		id := r.extractID(key)
		if id == "" {
			continue
		}
		if n, err := strconv.Atoi(id); err == nil && n > maxNum {
			maxNum = n
		}
	}

	return strconv.Itoa(maxNum + 1), nil
}

// extractID returns the id portion of a redis key. The underlying redis client
// prepends a global prefix (e.g. "soliant-mock-api:") to every key and SCAN
// returns the full key, so a simple TrimPrefix is not enough. This helper
// finds the repository prefix anywhere in the key and returns whatever follows.
func (r *Repository[T, PT]) extractID(key string) string {
	idx := strings.Index(key, r.prefix)
	if idx < 0 {
		return ""
	}
	return key[idx+len(r.prefix):]
}

// deepMerge recursively merges src into dst. Nested objects are merged;
// primitive values and arrays from src replace the corresponding value in dst.
func deepMerge(dst, src map[string]any) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				deepMerge(dstMap, srcMap)
				continue
			}
		}
		dst[key] = value
	}
}

