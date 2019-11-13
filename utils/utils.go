package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
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

func HexStringToBytes(s string) []byte {
	return BigIntToBytes(NewBigIntFromHexString(s))
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

func BigIntToBytes(x *big.Int) []byte {
	return x.Bytes()
}

func NewBigIntFromHexString(s string) *big.Int {
	x, err := new(big.Int).SetString(s, 16)
	if err == false {
		return nil
	}
	return x
}

// func NewBigIntFromBigInt(x *big.Int)

func NewBigIntFromBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func NewBigIntFromOne() *big.Int {
	return new(big.Int).SetInt64(1)
}

func NewBigIntFromZero() *big.Int {
	return new(big.Int).SetInt64(0)
}

func BigIntToHexString(x *big.Int) string {
	return fmt.Sprintf("%x", x.Bytes())
}

// BytesXor xs ^= ys
func BytesXor(xs []byte, ys []byte, commonLen int) {
	for i := 0; i != commonLen; i++ {
		xs[i] ^= ys[i]
	}
}

// BigIntAdd 大数加
func BigIntAdd(x, y *big.Int) *big.Int {
	var z big.Int
	z.Add(x, y)
	return &z
}

// BigIntSub 大数剪
func BigIntSub(x, y *big.Int) *big.Int {
	var z big.Int
	z.Sub(x, y)
	return &z
}

// BigIntMod 大数模
func BigIntMod(x, y *big.Int) *big.Int {
	var z big.Int
	z.Mod(x, y)
	return &z
}

// BigIntModInverse 大数模逆
func BigIntModInverse(x, y *big.Int) *big.Int {
	var z big.Int
	z.ModInverse(x, y)
	return &z
}

// BigIntMul 大数乘
func BigIntMul(x, y *big.Int) *big.Int {
	var z big.Int
	z.Mul(x, y)
	return &z
}

// BigIntLsh 大数左移
func BigIntLsh(x *big.Int, n uint) *big.Int {
	var z big.Int
	z.Lsh(x, n)
	return &z
}

// BigIntSetBit 大数设置比特位
func BigIntSetBit(x *big.Int, i int, b uint) *big.Int {
	var z big.Int
	z.SetBit(x, i, b)
	return &z
}

// BigIntAnd 大数与
func BigIntAnd(x, y *big.Int) *big.Int {
	var z big.Int
	z.And(x, y)
	return &z
}

// BitLenToBytesLen 比特流长度转字节流长度
func BitsLenToBytesLen(bitlen int) int {
	return (bitlen + 7) >> 3
}
