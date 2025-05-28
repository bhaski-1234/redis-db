package protocol

import (
	"errors"
	"strings"

	"github.com/bhaski-1234/redis-db/constant"
)

func DecodeInteger(data []byte) (int, int, error) {
	if len(data) < 2 || data[0] != ':' {
		return 0, 0, errors.New(constant.ErrInvalidRESP)
	}
	var value int
	i := 1
	sign := 1
	if i < len(data) && data[i] == '-' {
		sign = -1
		i++
	}
	for i < len(data) && data[i] != '\r' {
		if data[i] < '0' || data[i] > '9' {
			return 0, 0, errors.New(constant.ErrInvalidRESP)
		}
		value = value*10 + int(data[i]-'0')
		i++
	}
	return sign * value, i + 2, nil
}

func DecodeBulkString(data []byte) (string, int, error) {
	if len(data) < 2 || data[0] != '$' {
		return "", 0, errors.New(constant.ErrInvalidRESP)
	}

	var length int
	i := 1
	for i < len(data) && data[i] != '\r' {
		if data[i] < '0' || data[i] > '9' {
			return "", 0, errors.New(constant.ErrInvalidRESP)
		}
		length = length*10 + int(data[i]-'0')
		i++
	}

	i += 2 // skip \r\n
	return string(data[i : i+length]), i + length + 2, nil
}

func DecodeSimpleString(data []byte) (string, int, error) {
	if len(data) < 2 || data[0] != '+' {
		return "", 0, errors.New(constant.ErrInvalidRESP)
	}
	i := 1
	var builder strings.Builder
	for i < len(data) && data[i] != '\r' {
		builder.WriteByte(data[i])
		i++
	}
	return builder.String(), i + 2, nil
}

func DecodeArray(data []byte) ([]interface{}, int, error) {
	if len(data) < 2 || data[0] != '*' {
		return nil, 0, errors.New(constant.ErrInvalidRESP)
	}

	var length int
	i := 1
	for i < len(data) && data[i] != '\r' {
		if data[i] < '0' || data[i] > '9' {
			return nil, 0, errors.New(constant.ErrInvalidRESP)
		}
		length = length*10 + int(data[i]-'0')
		i++
	}

	i += 2 // skip \r\n
	result := make([]interface{}, length)
	for j := 0; j < length; j++ {
		item, nextIndex, err := DecodeRESP(data[i:])
		if err != nil {
			return nil, 0, err
		}
		result[j] = item
		i += nextIndex
	}
	return result, i, nil
}

func DecodeError(data []byte) (string, int, error) {
	if len(data) < 2 || data[0] != '-' {
		return "", 0, errors.New(constant.ErrInvalidRESP)
	}
	i := 1
	var builder strings.Builder
	for i < len(data) && data[i] != '\r' {
		builder.WriteByte(data[i])
		i++
	}
	return builder.String(), i + 2, nil
}

func DecodeRESP(data []byte) (interface{}, int, error) {
	switch data[0] {
	case ':':
		return DecodeInteger(data)
	case '$':
		return DecodeBulkString(data)
	case '+':
		return DecodeSimpleString(data)
	case '*':
		return DecodeArray(data)
	case '-':
		return DecodeError(data)
	default:
		return nil, 0, errors.New(constant.ErrInvalidRESP)
	}
}
