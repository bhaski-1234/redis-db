package command

import "github.com/bhaski-1234/redis-db/storage/inMemory"

func HandleGet(args []string) (interface{}, error) {
	inmemory := inMemory.GetInMemoryStore()
	key := args[1]
	value, exists := inmemory.Get(key)
	if !exists {
		return nil, nil // Key does not exist
	}
	return value, nil
}

func HandleSet(args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil // Not enough arguments
	}
	key := args[1]
	value := args[2]
	inmemory := inMemory.GetInMemoryStore()
	inmemory.Set(key, value)
	return "OK", nil
}
