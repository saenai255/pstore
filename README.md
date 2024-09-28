# Persistent Storage (pstore)

`pstore` is a Go library designed to provide simple and efficient persistent key-value storage with support for in-memory caching and optional thread-safety. You can store, retrieve, and manage data with minimal boilerplate, with automatic support for saving and loading to disk.

## Features

- **Persistent Storage**: Stores key-value pairs in memory and optionally saves them to disk.
- **Configurable Caching**: Control the number of items cached in memory or disable limits.
- **Thread-Safe Access**: Optionally use thread-safe access for concurrent operations.
- **In-Memory Only Option**: Create purely in-memory caches without persistence.
- **Flexible Storage Management**: Save individual key-value pairs or all items in a single file.
- **Error Type Identification**: Easily identify various error types like deletion, serialization, and disk read failures.

## Installation

```bash
go get github.com/saenai255/pstore
```

## Usage

### Basic Initialization

Create a new persistent storage cache:

```go
ps := pstore.New("/path/to/storage", "cache_name")
```

Or create an in-memory-only storage cache (no disk persistence):

```go
ps := pstore.NewInMemory("memory_cache")
```

### Configuring the Storage

Configure options after initialization to customize the behavior of `PersistentStorage`.

- **MaxMemItems**: Set the maximum number of items to be stored in memory (default: 100). Use `pstore.MEM_ITEMS_UNLIMITED` for unlimited items.
- **ThreadSafe**: Enable thread-safe access for concurrent operations.
- **SaveToDiskOnSet**: Control whether data is automatically saved to disk when a key is set (default: `true`).
- **SingleCacheFile**: Store all key-value pairs in a single file when saving to disk (default: `false`).

```go
ps := pstore.New("/path/to/storage", "cache_name")

// Configure options
ps.MaxMemItems = pstore.MemoryItemsCount(200) // Store up to 200 items in memory
ps.ThreadSafe = true                          // Enable thread-safe access
ps.SaveToDiskOnSet = false                    // Save to disk manually
ps.SingleCacheFile = true                     // Save all items to a single file
```

### Storing and Retrieving Data

Set a key-value pair:

```go
err := ps.Set("key1", "some_value")
if err != nil {
    fmt.Println("Error setting value:", err)
}
```

Get a key-value pair:

```go
var value string
err := ps.Get("key1", &value)
if err != nil {
    fmt.Println("Error getting value:", err)
} else {
    fmt.Println("Value:", value)
}
```

### Checking and Deleting Keys

Check if a key exists in the cache:

```go
exists := ps.Has("key1")
fmt.Println("Key exists:", exists)
```

Delete a key from the cache:

```go
err := ps.Delete("key1")
if err != nil {
    fmt.Println("Error deleting key:", err)
}
```

### Retrieving Cache Metadata

Get all keys in the cache:

```go
keys := ps.Keys()
fmt.Println("Cache keys:", keys)
```

Get the number of items in the cache:

```go
count := ps.Len()
fmt.Println("Cache item count:", count)
```

### Saving and Loading Data

Save the current state of the cache to disk:

```go
err := ps.SaveToDisk()
if err != nil {
    fmt.Println("Error saving to disk:", err)
}
```

### Handling Errors

The library provides utility functions to identify different error types:

```go
if pstore.IsDeleteFailed(err) {
    fmt.Println("Delete operation failed")
}

if pstore.IsSaveToDiskFailed(err) {
    fmt.Println("Failed to save to disk")
}

if pstore.IsKeyNotFound(err) {
    fmt.Println("Key not found")
}
```

## Examples

### Creating a Thread-Safe Persistent Storage with Unlimited Cache

```go
ps := pstore.New("/path/to/storage", "cache_name")

ps.MaxMemItems = pstore.MEM_ITEMS_UNLIMITED
ps.ThreadSafe = true
```

### Save All Cache Items to a Single File on Disk

```go
ps := pstore.New("/path/to/storage", "single_cache")

ps.SingleCacheFile = true
```

### Customizing Cache Behavior

Create a persistent storage with a maximum of 50 items in memory and custom disk-saving behavior:

```go
ps := pstore.New("/path/to/storage", "my_cache")

ps.MaxMemItems = pstore.MemoryItemsCount(50)
ps.SaveToDiskOnSet = false
ps.SingleCacheFile = true
```

## License

This library is licensed under the MIT License.

---

With `pstore`, you can easily create and manage key-value pairs, ensuring that your data is cached efficiently and safely persisted to disk when needed. Enjoy streamlined storage management for your Go applications!