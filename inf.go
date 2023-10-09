package partdb

type Instance[K comparable, V any] interface {
	// Index takes a side, a key, and a key slice. The key is added to the side if it does not exist.
	// If the keys in the slice do not exist on the opposite side, they are created.
	// All keys are two way indexed. Existing keys are modified.
	Index(bool, K, ...K) error

	// Remove takes a side and key, deleting the key from the side.
	// All edges to this key are removed before deletion of the key.
	Remove(bool, K) error

	// GetKey takes a side and a key, returning the value of the key, or returning an error if key does not exist.
	GetKey(bool, K) (*V, error)

	// LiseKeys takes a side and a key, returning the index map of that key.
	ListKey(bool, K) (map[K]V, error)

	// ListPart takes a side and returns a map of all keys and their values on that side.
	ListPart(bool) (map[K]V, error)
}

type PersistentInstance[K comparable, V any] interface {
	Instance[K, V]

	// Serializes the instance into a file at the path argument. Overwrites or creates file.
	Serialize(path string) error

	// Deserializes the instance from a file at the path argument.
	Deserialize(path string) error
}
