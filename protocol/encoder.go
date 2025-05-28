package protocol

import (
	"fmt"
	"log"
	"strconv"
)

func EncodeInteger(value int) []byte {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	return []byte(":" + sign + strconv.Itoa(value) + "\r\n")
}

func EncodeBulkString(value string) []byte {
	return []byte("$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n")
}

func EncodeSimpleString(value string) []byte {
	return []byte("+" + value + "\r\n")
}

func EncodeError(value string) []byte {
	return []byte("-" + value + "\r\n")
}

func EncodeArray(values []interface{}) []byte {
	if len(values) == 0 {
		return []byte("*0\r\n")
	}
	var result []byte
	result = append(result, []byte("*")...)
	result = append(result, []byte(strconv.Itoa(len(values))+"\r\n")...)
	for _, value := range values {
		switch v := value.(type) {
		case int:
			result = append(result, EncodeInteger(v)...)
		case string:
			result = append(result, EncodeBulkString(v)...)
		default:
			log.Printf("Warning: Unsupported type %T encountered in encodeArray", value)
		}
	}
	return result
}

func EncodeResponse(data interface{}) []byte {
	switch v := data.(type) {
	case string:
		return EncodeSimpleString(v)
	case int:
		return EncodeInteger(v)
	case []interface{}:
		return EncodeArray(v)
	case error:
		return EncodeError(v.Error())
	default:
		// Default to bulk string for other types
		return EncodeBulkString(fmt.Sprintf("%v", v))
	}
}
