package sm3

import (
	"bytes"
	"encoding/binary"
	"math/bits"
)

const SIZE = 32
const BLOCKSIZE = 64

// -----------------------------------------------------------------------------
// 数据结构
// -----------------------------------------------------------------------------

// Context 杂凑过程上下文
type Context struct {
	state    [8]uint32
	bitCount uint64 // GB/T 32905 5.1 l < 2**64
	buffer   [16]uint32
}

type Block [16]uint32

type LastBlock struct {
	content [14]uint32
	length  uint64
}

type PaddedBlock struct {
	blocks     []Block
	bitCount   int
	n          int
	lastBlocks [2]LastBlock
}

// -----------------------------------------------------------------------------
// GB/T 32905 3
// -----------------------------------------------------------------------------
var rotl32 = bits.RotateLeft32

// -----------------------------------------------------------------------------
// GB/T 32905 4
// -----------------------------------------------------------------------------

// iv 4.1
var iv = [8]uint32{
	0x7380166f,
	0x4914b2b9,
	0x172442d7,
	0xda8a0600,
	0xa96f30bc,
	0x163138aa,
	0xe38dee4d,
	0xb0fb0e4e}

// t 4.2
func t(j int) uint32 {
	if j >= 16 {
		return 0x7a879d8a
	} else {
		return 0x79cc4519
	}
}

// ff 4.3
func ff(j int, x, y, z uint32) uint32 {
	if j >= 16 {
		return ((x | y) & (x | z) & (y | z))
	} else {
		return x ^ y ^ z
	}
}

// gg 4.3
func gg(j int, x, y, z uint32) uint32 {
	if j >= 16 {
		return ((x & y) | ((^x) & z))
	} else {
		return x ^ y ^ z
	}
}

// p0 4.4
func p0(x uint32) uint32 {
	return x ^ rotl32(x, 9) ^ rotl32(x, 17)
}

// p1 4.4
func p1(x uint32) uint32 {
	return x ^ rotl32(x, 15) ^ rotl32(x, 23)
}

// -----------------------------------------------------------------------------
// GB/T 32905 5
// -----------------------------------------------------------------------------

// Init 5.1
func Init() *Context {
	var ctx Context
	for i := 0; i < 8; i++ {
		ctx.state[i] = iv[i]
	}
	return &ctx
}

// Padding 5.2 填充
func Padding(message []byte) []byte {
	msgLen := len(message)
	msgBuf := bytes.NewBuffer(message)
	msgBuf.WriteByte(0x80)
	for msgBuf.Len()%BLOCKSIZE != 56 {
		msgBuf.WriteByte(0x00)
	}
	lenBytes := new(bytes.Buffer)
	binary.Write(lenBytes, binary.BigEndian, uint64(msgLen*8))
	msgBuf.ReadFrom(lenBytes)
	return msgBuf.Bytes()
}

// MessageExpansion 5.3.2 消息扩展
func (ctx *Context) MessageExpansion(w *[68]uint32, wp *[64]uint32) {
	for i := 0; i < 16; i++ {
		w[i] = ctx.buffer[i]
	}

	for j := 16; j < 68; j++ {
		w[j] = p1(w[j-16]^w[j-9]^rotl32(w[j-3], 15)) ^ rotl32(w[j-13], 7) ^ w[j-6]
	}

	for j := 0; j < 64; j++ {
		wp[j] = w[j] ^ w[j+4]
	}
}

// Compress 5.3 迭代压缩
// TODO: 需要将 5.3.2, 5.3.3 和 5.4 分开
// TODO: 需要去除 n
func (ctx *Context) Compress(msg []byte, n int) {
	var (
		w      [68]uint32
		wprime [64]uint32
	)
	v := iv

	for len(msg) >= 64 {
		// 5.3.2 消息扩展

		// 5.3.2 a)
		for i := 0; i < 16; i++ {
			w[i] = binary.BigEndian.Uint32(msg[4*i : 4*(i+1)])
		}

		// 5.3.2 b)
		for j := 16; j < 68; j++ {
			w[j] = p1(w[j-16]^w[j-9]^rotl32(w[j-3], 15)) ^ rotl32(w[j-13], 7) ^ w[j-6]
		}

		// 5.3.2 c)
		for j := 0; j < 64; j++ {
			wprime[j] = w[j] ^ w[j+4]
		}

		// 5.3.3 压缩函数
		var ss1, ss2, tt1, tt2 uint32

		a, b, c, d, e, f, g, h :=
			v[0],
			v[1],
			v[2],
			v[3],
			v[4],
			v[5],
			v[6],
			v[7]

		for j := 0; j < 64; j++ {
			ss1 = rotl32(rotl32(a, 12)+e+rotl32(t(j), j), 7)
			ss2 = ss1 ^ rotl32(a, 12)
			tt1 = ff(j, a, b, c) + d + ss2 + wprime[j]
			tt2 = gg(j, e, f, g) + h + ss1 + w[j]
			d = c
			c = rotl32(b, 9)
			b = a
			a = tt1
			h = g
			g = rotl32(f, 19)
			f = e
			e = p0(tt2)
		}
		v[0] ^= a
		v[1] ^= b
		v[2] ^= c
		v[3] ^= d
		v[4] ^= e
		v[5] ^= f
		v[6] ^= g
		v[7] ^= h
		msg = msg[64:]
	}

	// 5.4 杂凑值
	for i, word := range v {
		ctx.state[i] = word
	}
}
