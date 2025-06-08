package command

import (
	diskstorage "github.com/bhaski-1234/redis-db/storage/diskStorage"
	"github.com/bhaski-1234/redis-db/storage/inMemory"
	"github.com/bhaski-1234/redis-db/utils"
)

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
	// Check if an expiration time is provided
	if len(args) > 4 && args[3] == "EX" {
		expiration, err := utils.ParseDuration(args[4])
		if err != nil {
			return nil, err // Invalid duration format
		}
		inmemory.SetWithExpiration(key, value, expiration)
	}
	return "OK", nil
}

func HandleDel(args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil
	}
	key := args[1]
	inmemory := inMemory.GetInMemoryStore()
	inmemory.Delete(key)
	inmemory.DeleteExpiration(key)
	return 1, nil
}

func HandleExists(args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil
	}
	key := args[1]
	inmemory := inMemory.GetInMemoryStore()
	exists := inmemory.Exists(key)
	if exists {
		return 1, nil
	}
	return 0, nil
}

func HandleTTL(args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, nil // Not enough arguments
	}
	key := args[1]
	inmemory := inMemory.GetInMemoryStore()
	ttl := inmemory.GetTTL(key)
	if ttl < 0 {
		return -1, nil // Key does not exist or has no expiration
	}
	return ttl, nil
}

func HandleSave(args []string) (interface{}, error) {
	disk := diskstorage.NewDiskStorage()
	if err := disk.Save("dump"); err != nil {
		return nil, err // Handle save error
	}
	return "OK", nil
}
