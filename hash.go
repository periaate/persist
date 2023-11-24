package persist

import (
	"github.com/OneOfOne/xxhash"
)

// HashBytes is an alias for xxhash.Checksum64.
var HashBytes = xxhash.Checksum64

// HashString is an alias for xxhash.ChecksumString64.
var HashString = xxhash.ChecksumString64

const (
	prime uint64 = 2654435761
)

// NewHashU64 returns a new instance of a Murmur3 function. The returned function is not thread safe.
func NewHashU64() func(uint64) uint64 {
	var hash uint64
	return func(key uint64) uint64 {
		hash = (prime + (key >> 0)) * prime
		hash = (hash << 13) | (hash >> 19)
		hash *= prime

		hash += (key >> 8) * prime
		hash = (hash << 13) | (hash >> 19)
		hash *= prime

		hash += (key >> 16) * prime
		hash = (hash << 13) | (hash >> 19)
		hash *= prime

		hash += (key >> 24) * prime
		hash = (hash << 13) | (hash >> 19)
		hash *= prime

		hash ^= 4
		hash ^= hash >> 16
		hash *= prime
		hash ^= hash >> 13
		hash *= prime
		hash ^= hash >> 16

		return hash
	}
}

// HashU64 is an implementation of a Murmur3 algorithm, manually unwrapped.
func HashU64(key uint64) (hash uint64) {
	hash = (prime + (key >> 0)) * prime
	hash = (hash << 13) | (hash >> 19)
	hash *= prime

	hash += (key >> 8) * prime
	hash = (hash << 13) | (hash >> 19)
	hash *= prime

	hash += (key >> 16) * prime
	hash = (hash << 13) | (hash >> 19)
	hash *= prime

	hash += (key >> 24) * prime
	hash = (hash << 13) | (hash >> 19)
	hash *= prime

	hash ^= 4
	hash ^= hash >> 16
	hash *= prime
	hash ^= hash >> 13
	hash *= prime
	hash ^= hash >> 16

	return hash
}
