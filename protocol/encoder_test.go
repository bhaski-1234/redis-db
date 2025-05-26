package protocol

import "testing"

func TestEncodedBulkString(t *testing.T) {
	input := []string{
		"Hello",
		"",
		"Hello World",
	}

	output := [][]byte{
		[]byte("$5\r\nHello\r\n"),
		[]byte("$0\r\n\r\n"),
		[]byte("$11\r\nHello World\r\n"),
	}

	for i, data := range input {
		result := encodeBulkString(data)
		if string(result) != string(output[i]) {
			t.Errorf("TestEncodedBulkString failed for input %s: expected %s, got %s", data, output[i], result)
		}
	}
}

func TestEncodedSimpleString(t *testing.T) {
	input := []string{
		"OK",
		"Hello World",
	}

	output := [][]byte{
		[]byte("+OK\r\n"),
		[]byte("+Hello World\r\n"),
	}

	for i, data := range input {
		result := encodeSimpleString(data)
		if string(result) != string(output[i]) {
			t.Errorf("TestEncodedSimpleString failed for input %s: expected %s, got %s", data, output[i], result)
		}
	}
}

func TestEncodedInteger(t *testing.T) {
	input := []int{
		12345,
		0,
		-67890,
	}

	output := [][]byte{
		[]byte(":12345\r\n"),
		[]byte(":0\r\n"),
		[]byte(":-67890\r\n"),
	}

	for i, data := range input {
		result := encodeInteger(data)
		if string(result) != string(output[i]) {
			t.Errorf("TestEncodedInteger failed for input %d: expected %s, got %s", data, output[i], result)
		}
	}
}

func TestEncodedArray(t *testing.T) {
	input := [][]interface{}{
		{"foo", "Hello", "bar"},
		{7, 42},
		{"empty", 13, "world"},
		{},
	}

	output := [][]byte{
		[]byte("*3\r\n$3\r\nfoo\r\n$5\r\nHello\r\n$3\r\nbar\r\n"),
		[]byte("*2\r\n:7\r\n:42\r\n"),
		[]byte("*3\r\n$5\r\nempty\r\n:13\r\n$5\r\nworld\r\n"),
		[]byte("*0\r\n"),
	}

	for i, data := range input {
		result := encodeArray(data)
		if string(result) != string(output[i]) {
			t.Errorf("TestEncodedArray failed for input %v: expected %s, got %s", data, output[i], result)
		}
	}
}
