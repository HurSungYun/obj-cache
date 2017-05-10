package objcache

import "time"

// Config for cache
// TODO: memory limit or entry limit
type Config struct {
	MaxEntryLimit int
	Expiration    time.Duration
}
