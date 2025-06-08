package diskstorage

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/bhaski-1234/redis-db/constant"
	"github.com/bhaski-1234/redis-db/storage/inMemory"
	"github.com/bhaski-1234/redis-db/utils"
)

type DiskStorage struct {
	inMemoryStore *inMemory.InMemoryStore
}

func NewDiskStorage() *DiskStorage {
	return &DiskStorage{
		inMemoryStore: inMemory.GetInMemoryStore(),
	}
}

func (ds *DiskStorage) Save(fileName string) error {
	//format
	// [Header] [Type][KeyLen][Key][ValLen][Val] [Type][KeyLen][Key][ValLen][Val] [EOF]
	// Implement the logic to save data to a file
	fs, err := os.OpenFile(fileName+constant.DataFileExtension, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()
	// Write header
	fs.Write([]byte(constant.Header))

	// Save all key-values
	ds.inMemoryStore.Store.Range(func(key, value interface{}) bool {
		keyStr, _ := key.(string)
		switch value.(type) {
		case string:
			writeString(fs, keyStr, value.(string))
		case int:
			writeInt(fs, keyStr, value.(int))
		case []byte:
			//TODO
		}
		return true
	})

	// Save all TTL information
	ds.inMemoryStore.GetExpirations(func(key string, expTime time.Time) bool {
		writeTTL(fs, key, expTime)
		return true
	})

	fs.Write([]byte{constant.EOF})
	return nil
}

func (ds *DiskStorage) Load(fileName string) error {
	ds.inMemoryStore.Clear()
	fs, err := os.OpenFile(fileName+constant.DataFileExtension, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()

	// Read and verify header
	header := make([]byte, len(constant.Header))
	if _, err := fs.Read(header); err != nil {
		return err
	}
	if string(header) != constant.Header {
		return errors.New("invalid file header")
	}

	for {
		// Read type
		typeBuf := make([]byte, 1)
		if _, err := fs.Read(typeBuf); err != nil {
			return err
		}
		if typeBuf[0] == constant.EOF {
			break
		}

		// Read key length (varint)
		keyLen, err := readVarInt(fs)
		if err != nil {
			return err
		}
		// Read key
		keyBuf := make([]byte, keyLen)
		if _, err := fs.Read(keyBuf); err != nil {
			return err
		}

		// Read value length (varint)
		valLen, err := readVarInt(fs)
		if err != nil {
			return err
		}
		// Read value
		valBuf := make([]byte, valLen)
		if _, err := fs.Read(valBuf); err != nil {
			return err
		}

		keyStr := string(keyBuf)
		switch typeBuf[0] {
		case constant.TypeString:
			ds.inMemoryStore.Set(keyStr, string(valBuf))
		case constant.TypeInteger:
			intVal, _ := strconv.Atoi(string(valBuf))
			ds.inMemoryStore.Store.Store(keyStr, intVal)
		case constant.TypeTTL:
			// TTL is stored as timestamp in milliseconds
			ttlMs, _ := strconv.ParseInt(string(valBuf), 10, 64)
			expTime := time.Unix(0, ttlMs*int64(time.Millisecond))

			// Skip expired keys
			if time.Now().Before(expTime) {
				ds.inMemoryStore.SetExpiration(keyStr, expTime)
			} else {
				// If expired, delete the key
				ds.inMemoryStore.Delete(keyStr)
			}
		}
	}
	return nil
}

// Helper to read a varint from file
func readVarInt(fs *os.File) (int, error) {
	var buf []byte
	for {
		b := make([]byte, 1)
		if _, err := fs.Read(b); err != nil {
			return 0, err
		}
		buf = append(buf, b[0])
		if b[0]&0x80 == 0 {
			break
		}
	}
	return int(utils.DecodeVarIntBigEndian(buf)), nil
}

func writeString(fs *os.File, key string, value string) error {
	fs.Write([]byte{constant.TypeString})
	// Get the varint length of the key and value
	keyBytes := utils.EncodeVarIntBigEndian(len(key))
	valueBytes := utils.EncodeVarIntBigEndian(len(value))
	fs.Write(keyBytes)
	fs.Write([]byte(key))
	fs.Write(valueBytes)
	fs.Write([]byte(value))
	return nil
}

func writeInt(fs *os.File, key string, value int) error {
	valueStr := strconv.Itoa(value)
	fs.Write([]byte{constant.TypeInteger})
	// Get the varint length of the key
	keyBytes := utils.EncodeVarIntBigEndian(len(key))
	valueBytes := utils.EncodeVarIntBigEndian(len(valueStr))
	fs.Write(keyBytes)
	fs.Write([]byte(key))
	fs.Write(valueBytes)
	fs.Write([]byte(valueStr))
	return nil
}

func writeTTL(fs *os.File, key string, expTime time.Time) error {
	fs.Write([]byte{constant.TypeTTL})
	// Write key
	keyBytes := utils.EncodeVarIntBigEndian(len(key))
	fs.Write(keyBytes)
	fs.Write([]byte(key))

	// Store TTL as timestamp in milliseconds
	ttlMs := expTime.UnixNano() / int64(time.Millisecond)
	ttlStr := strconv.FormatInt(ttlMs, 10)

	// Write TTL value
	valueBytes := utils.EncodeVarIntBigEndian(len(ttlStr))
	fs.Write(valueBytes)
	fs.Write([]byte(ttlStr))
	return nil
}
