package utils

import (
	"bytes"
	"encoding/binary"
)

// WordsToBytes 字slice转化成网络流
func WordsToBytes(words []uint32) []byte {
	buf := new(bytes.Buffer)
	for _, word := range words {
		binary.Write(buf, binary.BigEndian, word)
	}
	return buf.Bytes()
}

func BytesTo16Words(stream []byte) [16]uint32 {
	buf := bytes.NewBuffer(stream)
	var ans [16]uint32
	binary.Read(buf, binary.BigEndian, &ans)
	return ans
}

func BytesTo14Words(stream []byte) [14]uint32 {
	buf := bytes.NewBuffer(stream)
	var ans [14]uint32
	binary.Read(buf, binary.BigEndian, &ans)
	return ans
}
