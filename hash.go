package partdb

import (
	"encoding/binary"
	"hash/fnv"
)

func Fnv_u64() func(uint64) uint64 {
	f := fnv.New64()
	return func(key uint64) uint64 {
		fnvBytes := make([]byte, 8)
		f.Reset()
		binary.LittleEndian.PutUint64(fnvBytes, key)
		f.Write(fnvBytes)
		return f.Sum64()
	}
}

func Fnv_u32() func(uint32) uint32 {
	f := fnv.New32()
	return func(key uint32) uint32 {
		fnvBytes := make([]byte, 4)
		f.Reset()
		binary.LittleEndian.PutUint32(fnvBytes, key)
		f.Write(fnvBytes)
		return f.Sum32()
	}
}

func Hash_string() func(string) uint64 {
	h := Hash_bytes()
	return func(key string) uint64 {
		return h([]byte(key))
	}
}

func Hash_bytes32() func([32]byte) uint64 {
	h := Hash_bytes()
	return func(key [32]byte) uint64 {
		return h(key[:])
	}
}

func Hash_bytes() func([]byte) uint64 {
	hfn := Hash_u64()
	return func(key []byte) uint64 {
		keyN := binary.LittleEndian.Uint64(key[:8]) // huh?
		return hfn(keyN)
	}
}

// Creates a murmur3 hash function with its own allocated values
func Hash_u64() func(uint64) uint64 {
	var hash uint64
	const (
		prime uint64 = 2654435761
	)

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
