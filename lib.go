package pstore

import (
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
)

type any = interface{}

type memItemCount int

const MEM_ITEMS_UNLIMITED memItemCount = -1
const MEM_ITEMS_DEFAULT memItemCount = 100

func MemoryItemsCount(count int) memItemCount {
	if count < 0 {
		return MEM_ITEMS_UNLIMITED
	}

	return memItemCount(count)
}

type PersistentStorage struct {
	// The maximum number of items to keep in memory.
	MaxMemItems memItemCount
	// If true, the cache will be thread-safe.
	ThreadSafe bool
	// If true, the cache will be saved to disk when a key is set. Default is true.
	SaveToDiskOnSet bool
	// If true, the cache will be saved to disk when the cache is saved. Default is false.
	SingleCacheFile bool

	path     string
	name     string
	cache    map[string]any
	inMemory bool
	mutex    *sync.Mutex
}

// New creates a new PersistentStorage instance.
//
// Parameters:
//   - path: The path to the directory where the cache files will be stored.
//   - name: The name of the cache.
//
// Returns:
//   - A new PersistentStorage instance.
func New(path, name string) *PersistentStorage {
	return &PersistentStorage{
		path:            path,
		name:            name,
		cache:           make(map[string]any),
		MaxMemItems:     MEM_ITEMS_DEFAULT,
		inMemory:        false,
		SaveToDiskOnSet: true,
		ThreadSafe:      false,
		mutex:           new(sync.Mutex),
	}
}

// NewInMemory creates a new PersistentStorage instance that only stores data in memory.
//
// Parameters:
//   - name: The name of the cache.
//
// Returns:
//   - A new PersistentStorage instance.
func NewInMemory(name string) *PersistentStorage {
	return &PersistentStorage{
		name:            name,
		cache:           make(map[string]any),
		MaxMemItems:     MEM_ITEMS_DEFAULT,
		inMemory:        true,
		SaveToDiskOnSet: true,
		ThreadSafe:      false,
		mutex:           new(sync.Mutex),
	}
}

// Returns the number of items in the cache. Only the in-memory cache is counted.
//
// Returns:
//   - The number of items in the cache.
func (ps *PersistentStorage) Len() int {
	return len(ps.cache)
}

const error_delete_failed = "failed to delete"

// IsDeleteFailed returns true if the error is a delete failure.
func IsDeleteFailed(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_delete_failed)
}

// Delete removes the key from the cache and deletes the file from disk.
//
// Parameters:
//   - key: The key to delete.
//
// Returns:
//   - An error if the key does not exist or if the file could not be deleted.
func (ps *PersistentStorage) Delete(key string) error {
	delete(ps.cache, key)

	if !ps.inMemory {
		if err := os.RemoveAll(ps.getCachePath(key)); err != nil {
			return ps.errorf("%s %v: %v", error_delete_failed, key, err)
		}
	}

	return nil
}

// Has returns true if the key exists in the cache.
//
// Parameters:
//   - key: The key to check.
//
// Returns:
//   - True if the key exists in the cache.
func (ps *PersistentStorage) Has(key string) bool {
	if ps.ThreadSafe {
		ps.mutex.Lock()
		defer ps.mutex.Unlock()
	}

	_, ok := ps.cache[key]
	return ok
}

// Keys returns a list of all keys in the cache. Only the in-memory cache is counted. The keys are not sorted and the order is not guaranteed.
//
// Returns:
//   - A list of all keys in the cache.
func (ps *PersistentStorage) Keys() []string {
	if ps.ThreadSafe {
		ps.mutex.Lock()
		defer ps.mutex.Unlock()
	}

	keys := make([]string, 0, len(ps.cache))

	for k := range ps.cache {
		keys = append(keys, k)
	}

	return keys
}

// SaveToDisk saves the cache to disk. If SingleCacheFile is true, all keys are saved to a single file. Otherwise, each key is saved to a separate file. If SaveToDiskOnSet is true, the cache is saved to disk when a key is set and this method does nothing.
//
// Returns:
//   - An error if the cache could not be saved to disk.
func (ps *PersistentStorage) SaveToDisk() error {
	if ps.SaveToDiskOnSet {
		return nil
	}

	if ps.ThreadSafe {
		ps.mutex.Lock()
		defer ps.mutex.Lock()
	}

	if ps.SingleCacheFile {
		return ps.saveToDisk(single_cache_filename, nil)
	}

	for k, v := range ps.cache {
		if err := ps.saveToDisk(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (ps *PersistentStorage) set(key string, value any) error {
	ps.cache[key] = value
	if ps.SaveToDiskOnSet {
		if err := ps.saveToDisk(key, value); err != nil {
			return err
		}
	}

	if ps.MaxMemItems != MEM_ITEMS_UNLIMITED && len(ps.cache) > int(ps.MaxMemItems) {
		for k := range ps.cache {
			delete(ps.cache, k)
			break
		}
	}

	return nil
}

const cache_ext = ".pcache"

func (ps *PersistentStorage) getCachePath(key string) string {
	return path.Join(ps.path, ps.name+"_"+key+cache_ext)
}

const error_save_to_disk_failed = "failed to save"

// IsSaveToDiskFailed returns true if the error is a save failure.
func IsSaveToDiskFailed(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_save_to_disk_failed)
}

const error_serialize_failed = "failed to serialize"

// IsSerializeFailed returns true if the error is a serialization failure.
func IsSerializeFailed(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_serialize_failed)
}

const single_cache_filename = "single_full_cache"

func (ps *PersistentStorage) saveToDisk(key string, value any) error {
	if ps.inMemory {
		return nil
	}

	if err := os.MkdirAll(ps.path, 0755); err != nil {
		return ps.errorf("%s %v: %v", error_save_to_disk_failed, key, err)
	}

	if ps.SingleCacheFile {
		cacheBytes, err := serialize(ps.cache)
		if err != nil {
			return ps.errorf("%s %v: %v", error_serialize_failed, single_cache_filename, err)
		}

		if err := os.WriteFile(ps.getCachePath(single_cache_filename), cacheBytes, 0644); err != nil {
			return ps.errorf("%s %v: %v", error_save_to_disk_failed, single_cache_filename, err)
		}

		return nil
	}

	bytes, err := serialize(value)
	if err != nil {
		return ps.errorf("%s %v: %v", error_serialize_failed, key, err)
	}

	if err := os.WriteFile(ps.getCachePath(key), bytes, 0644); err != nil {
		return ps.errorf("%s %v: %v", error_save_to_disk_failed, key, err)
	}

	return nil
}

const error_expected_pointer = "expected pointer"

// IsExpectedPointer returns true if the error is an expected pointer error.
func IsExpectedPointer(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_expected_pointer)
}

func (ps *PersistentStorage) get(out any, key string) error {
	outReflect := reflect.ValueOf(out)
	if outReflect.Kind() != reflect.Ptr {
		return ps.errorf("%s but got type %v", error_expected_pointer, outReflect.Type())
	}

	it, ok := ps.cache[key]

	// If the key is not in the cache, read it from disk
	if !ok {
		err := ps.readFromDisk(out, key)

		// If the key is found on disk, cache it
		if err != nil {
			ps.cache[key] = outReflect.Elem().Interface()
		}

		return err
	}

	outReflect.Elem().Set(reflect.ValueOf(it))

	return nil
}

const error_key_not_found = "key not found"

// IsKeyNotFound returns true if the error is a key not found error.
func IsKeyNotFound(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_key_not_found)
}

const error_read_from_disk_failed = "failed to read from disk"

// IsReadFromDiskFailed returns true if the error is a read from disk failure.
func IsReadFromDiskFailed(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_read_from_disk_failed)
}

const error_deserialize_failed = "failed to deserialize"

// IsDeserializeFailed returns true if the error is a deserialization failure.
func IsDeserializeFailed(err error) bool {
	return IsPStoreError(err) && strings.Contains(err.Error(), error_deserialize_failed)
}

func (ps *PersistentStorage) readFromDisk(out interface{}, key string) error {
	if ps.inMemory {
		return ps.errorf("%s %v", error_key_not_found, key)
	}

	bytes, err := os.ReadFile(ps.getCachePath(key))
	if err != nil {
		if os.IsNotExist(err) {
			return ps.errorf("%s %v", error_key_not_found, key)
		}

		return ps.errorf("%s %v: %v", error_read_from_disk_failed, key, err)
	}

	if err := deserialize(bytes, out); err != nil {
		return ps.errorf("%s %v: %v", error_deserialize_failed, key, err)
	}

	return nil
}

// Set sets the value of the key in the cache. If the key already exists, it is overwritten. If the cache is thread-safe, the operation is atomic.
//
// Parameters:
//   - key: The key to set.
//   - value: The value to set. Should not be a pointer.
//
// Returns:
//   - An error if the value could not be set.
func (ps *PersistentStorage) Set(key string, value any) error {
	if ps.ThreadSafe {
		ps.mutex.Lock()
		defer ps.mutex.Lock()
	}

	var it any = value

	for reflect.TypeOf(it).Kind() == reflect.Ptr {
		it = reflect.ValueOf(it).Elem().Interface()
	}

	return ps.set(key, it)
}

// Get gets the value of the key from the cache. If the key does not exist in the cache, it is read from disk. If the key does not exist on disk, an error is returned. If the cache is thread-safe, the operation is atomic.
//
// Parameters:
//   - key: The key to get.
//   - out: The value to get. Must be a pointer.
//
// Returns:
//   - An error if the value could not be retrieved.
func (ps *PersistentStorage) Get(key string, out any) error {
	if ps.ThreadSafe {
		ps.mutex.Lock()
		defer ps.mutex.Lock()
	}

	return ps.get(out, key)
}
