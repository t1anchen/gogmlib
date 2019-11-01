package utils

import (
	"fmt"
	"testing"
)

func TestWordToBytes(t *testing.T) {
	input := []uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600,
		0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e}
	t.Logf("TestWordToBytes:input = %x", input)

	actual := WordsToBytes(input)
	actualStr := fmt.Sprintf("%x", actual)
	t.Logf("TestWordToBytes:actual = %x", actual)

	expectedStr := "7380166f4914b2b9172442d7da8a0600a96f30bc163138aae38dee4db0fb0e4e"
	t.Logf("TestWordToBytes:expectedStr = %s", expectedStr)
	if actualStr != expectedStr {
		t.Errorf(`utils:TestWordToBytes
expected=%s
actual=%s`, expectedStr, actualStr)
	}
}

func TestBytesTo16Words(t *testing.T) {
	var input []byte
	inputStr := "7380166f4914b2b9172442d7da8a0600a96f30bc163138aae38dee4db0fb0e4e7380166f4914b2b9172442d7da8a0600a96f30bc163138aae38dee4db0fb0e4e"
	fmt.Sscanf(inputStr, "%x", &input)
	t.Logf("TestBytesTo16Words:input = %x", input)

	actual := BytesTo16Words(input[:])
	actualStr := fmt.Sprintf("%x", actual)
	t.Logf("TestBytesTo16Words:actual = %x", actual)

	expected := []uint32{
		0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600,
		0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e,
		0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600,
		0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e}
	expectedStr := fmt.Sprintf("%x", expected)
	t.Logf("TestBytesTo16Words:expected = %x", actual)

	if actualStr != expectedStr {
		t.Errorf(`utils:TestBytesTo16Words
expected=%s
actual=%s`, expectedStr, actualStr)
	}

}
