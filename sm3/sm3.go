package sm3

import (
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math/bits"
)

const SIZE = 32
const BLOCKSIZE = 64

// -----------------------------------------------------------------------------
// 应用接口
// -----------------------------------------------------------------------------

// Checksum sm3面向流的 API
func Checksum(stream []byte) []byte {
	var ctx Context
	ctx.Reset()
	ctx.Write(stream)
	return ctx.Sum(nil)
}

// -----------------------------------------------------------------------------
// golang/hash/Hash 接口
// -----------------------------------------------------------------------------

// Context 杂凑过程上下文
type Context struct {
	hash     [8]uint32
	bitCount uint64 // GB/T 32905 5.1 l < 2**64
	buffer   [16]uint32
	m        []byte
}

// New 生成杂凑上下文实例
func New() hash.Hash {
	c := new(Context)
	c.Reset()
	return c
}

// Reset 实现 Hash 接口中的 Reset 函数，初始化杂凑上下文
func (ctx *Context) Reset() {
	for i, word := range iv {
		ctx.hash[i] = word
	}
	ctx.bitCount = 0
	ctx.buffer = [16]uint32{}
	ctx.m = []byte{}
}

// Sum 实现 Hash 接口中的 Sum 函数
func (ctx *Context) Sum(inputStream []byte) []byte {
	ctx.Write(inputStream)
	msg, err := ctx.Padding()
	if err != nil {
		log.Fatalf("sm3:Sum:Padding 失败, 错误为%v", err)
		return nil
	}

	// 最后一个block处理
	ctx.Compress(msg, len(msg)/ctx.BlockSize())

	// TODO: 需要使用 bytes 包和 binary 包的相关组件重构
	requiredLen := ctx.Size()
	if cap(inputStream)-len(inputStream) < requiredLen {
		newInputStream := make([]byte, len(inputStream), len(inputStream)+requiredLen)
		copy(newInputStream, inputStream)
		inputStream = newInputStream
	}

	outputStream := inputStream[len(inputStream) : len(inputStream)+requiredLen]

	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(outputStream[i*4:], ctx.hash[i])
	}
	return outputStream

}

// Size 实现 Hash 接口中的 Size 函数
func (ctx *Context) Size() int {
	return SIZE
}

// BlockSize 实现 Hash 接口中的 BlockSize 函数
func (ctx *Context) BlockSize() int {
	return BLOCKSIZE
}

// Write 实现 Hash 接口中 io.Writer 接口的 Write 函数
func (ctx *Context) Write(newChunk []byte) (count int, err error) {
	newChunkLen := len(newChunk)
	ctx.bitCount += uint64(len(newChunk) * 8)

	msg := append(ctx.m, newChunk...)
	n := len(msg) / ctx.BlockSize() // blocks
	ctx.Compress(msg, n)

	ctx.m = msg[n*ctx.BlockSize():]

	return newChunkLen, nil
}

// -----------------------------------------------------------------------------
// GB/T 32905 3
// -----------------------------------------------------------------------------
var rol32 = bits.RotateLeft32

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
	return x ^ rol32(x, 9) ^ rol32(x, 17)
}

// p1 4.4
func p1(x uint32) uint32 {
	return x ^ rol32(x, 15) ^ rol32(x, 23)
}

// -----------------------------------------------------------------------------
// GB/T 32905 5
// -----------------------------------------------------------------------------

// Padding 5.2 填充
func (ctx *Context) Padding() ([]byte, error) {

	msg := ctx.m
	msg = append(msg, 0x80)

	// l + 1 + k ≡ 448 (mod 512)
	for len(msg)%BLOCKSIZE != 56 {
		msg = append(msg, 0x00)
	}

	msg = append(msg, uint8(ctx.bitCount>>56&0xff))
	msg = append(msg, uint8(ctx.bitCount>>48&0xff))
	msg = append(msg, uint8(ctx.bitCount>>40&0xff))
	msg = append(msg, uint8(ctx.bitCount>>32&0xff))
	msg = append(msg, uint8(ctx.bitCount>>24&0xff))
	msg = append(msg, uint8(ctx.bitCount>>16&0xff))
	msg = append(msg, uint8(ctx.bitCount>>8&0xff))
	msg = append(msg, uint8(ctx.bitCount>>0&0xff))

	// 测试是否是 512 的倍数
	if len(msg)%64 != 0 {
		return nil, fmt.Errorf(`sm3:Context:Padding: 消息长度 = %d 不是 64字节/512位 的倍数`, len(msg))
	}
	return msg, nil
}

func (ctx *Context) computeSum() [SIZE]byte {
	return [SIZE]byte{0}
}

// MessageExpansion 5.3.2 消息扩展
func (ctx *Context) MessageExpansion(w *[68]uint32, wp *[64]uint32) {
	for i := 0; i < 16; i++ {
		w[i] = ctx.buffer[i]
	}

	for j := 16; j < 68; j++ {
		w[j] = p1(w[j-16]^w[j-9]^rol32(w[j-3], 15)) ^ rol32(w[j-13], 7) ^ w[j-6]
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
			w[j] = p1(w[j-16]^w[j-9]^rol32(w[j-3], 15)) ^ rol32(w[j-13], 7) ^ w[j-6]
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
			ss1 = rol32(rol32(a, 12)+e+rol32(t(j), j), 7)
			ss2 = ss1 ^ rol32(a, 12)
			tt1 = ff(j, a, b, c) + d + ss2 + wprime[j]
			tt2 = gg(j, e, f, g) + h + ss1 + w[j]
			d = c
			c = rol32(b, 9)
			b = a
			a = tt1
			h = g
			g = rol32(f, 19)
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
		ctx.hash[i] = word
	}
}
