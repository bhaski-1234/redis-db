package protocol

import "strconv"

func encodeInteger(value int) []byte {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	return []byte(":" + sign + strconv.Itoa(value) + "\r\n")
}

func encodeBulkString(value string) []byte {
	return []byte("$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n")
}

func encodeSimpleString(value string) []byte {
	return []byte("+" + value + "\r\n")
}

func encodeError(value string) []byte {
	return []byte("-" + value + "\r\n")
}

func encodeArray(values []interface{}) []byte {
	if len(values) == 0 {
		return []byte("*0\r\n")
	}
	var result []byte
	result = append(result, []byte("*")...)
	result = append(result, []byte(strconv.Itoa(len(values))+"\r\n")...)
	for _, value := range values {
		switch v := value.(type) {
		case int:
			result = append(result, encodeInteger(v)...)
		case string:
			result = append(result, encodeBulkString(v)...)
		default:
			log.Printf("Warning: Unsupported type %T encountered in encodeArray", value)
		}
	}
	return result
}
