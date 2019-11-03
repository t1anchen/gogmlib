package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
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

func BytesConcat(a ...[]byte) []byte {
	buf := new(bytes.Buffer)
	for _, x := range a {
		buf.Write(x)
	}
	return buf.Bytes()
}

func ReadBytesFromFileToBuffer(f *os.File) *bytes.Buffer {
	msgBuf := new(bytes.Buffer)
	block := make([]byte, 4096)
	r := bufio.NewReader(f)
	loaded, err := r.Read(block)
	for err == nil {
		msgBuf.Write(block[:loaded])
		loaded, err = r.Read(block)
	}
	return msgBuf
}

func ReadBytesFromStdinToBuffer() *bytes.Buffer {
	return ReadBytesFromFileToBuffer(os.Stdin)
}
