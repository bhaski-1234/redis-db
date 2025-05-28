package protocol

import (
	"testing"
)

func TestDecodedSimpleString(t *testing.T) {
	input := [][]byte{
		[]byte("+OK\r\n"),
		[]byte("+Hello World\r\n"),
	}

	output := [][]byte{
		[]byte("OK"),
		[]byte("Hello World"),
	}

	for i, data := range input {
		result, _, _ := DecodeRESP(data)
		resultConv, _ := result.(string)
		if resultConv != string(output[i]) {
			t.Errorf("TestSimpleString failed for input %s: expected %s, got %s", data, output[i], resultConv)
		}
	}
}

func TestDecodedInteger(t *testing.T) {
	input := [][]byte{
		[]byte(":12345\r\n"),
		[]byte(":0\r\n"),
		[]byte(":-67890\r\n"),
	}

	output := []int{
		12345,
		0,
		-67890,
	}

	for i, data := range input {
		result, _, _ := DecodeRESP(data)
		resultConv, _ := result.(int)
		if resultConv != output[i] {
			t.Errorf("TestInteger failed for input %s: expected %d, got %d", data, output[i], resultConv)
		}
	}
}

func TestDecodedBulkString(t *testing.T) {
	input := [][]byte{
		[]byte("$5\r\nHello\r\n"),
		[]byte("$0\r\n\r\n"),
		[]byte("$11\r\nHello World\r\n"),
		[]byte("$-1\r\n"), // nil bulk string
	}

	output := []string{
		"Hello",
		"",
		"Hello World",
		"", // nil bulk string is represented as an empty string
	}

	for i, data := range input {
		result, _, _ := DecodeRESP(data)
		resultConv, _ := result.(string)
		if resultConv != output[i] {
			t.Errorf("TestBulkString failed for input %s: expected %s, got %s", data, output[i], resultConv)
		}
	}
}

func TestDecodedArray(t *testing.T) {
	input := [][]byte{
		[]byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
		[]byte("*3\r\n:1\r\n:2\r\n:3\r\n"),
		[]byte("*0\r\n"),
	}

	output := [][]interface{}{
		{"foo", "bar"},
		{1, 2, 3},
		{},
	}

	for i, data := range input {
		result, _, _ := DecodeRESP(data)
		resultConv, _ := result.([]interface{})
		if len(resultConv) != len(output[i]) {
			t.Errorf("TestArray failed for input %s: expected length %d, got %d", data, len(output[i]), len(resultConv))
			continue
		}
		for j, item := range resultConv {
			if item != output[i][j] {
				t.Errorf("TestArray failed for input %s: expected %v, got %v", data, output[i][j], item)
			}
		}
	}
}

func TestDecodedErrorHandling(t *testing.T) {
	input := [][]byte{
		[]byte(":NotAnInteger\r\n"),
		[]byte("*InvalidArray\r\n"),
	}

	for _, data := range input {
		_, _, err := DecodeRESP(data)
		if err == nil {
			t.Errorf("Expected error for input %s, but got none", data)
		}
	}
}
