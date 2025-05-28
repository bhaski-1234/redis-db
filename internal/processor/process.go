package processor

import (
	"fmt"

	"github.com/bhaski-1234/redis-db/internal/dispatcher"
	"github.com/bhaski-1234/redis-db/protocol"
)

func Process(data []byte) (interface{}, error) {
	decoded, _, err := protocol.DecodeRESP(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RESP: %w", err)
	}

	// Check if decoded is a slice of interfaces
	decodedData, ok := decoded.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid decoded data type")
	}

	if len(decodedData) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd, ok := decodedData[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	d := dispatcher.NewDispatcher()
	args := make([]string, len(decodedData))

	for i, arg := range decodedData {
		if strArg, ok := arg.(string); ok {
			args[i] = strArg
		} else {
			// do nothing or handle error
			return nil, fmt.Errorf("invalid argument at position %d", i)
		}
	}

	result, err := d.Execute(cmd, args)
	return result, err
}
