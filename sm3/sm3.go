package sm3

import (
	"hash"
	"math/bits"
)

var (
	// GB/T 32907 4.1 初始值
	iv = [8]uint32{
		0x7380166f,
		0x4914b2b9,
		0x172442d7,
		0xda8a0600,
		0xa96f30bc,
		0x163138aa,
		0xe38dee4d,
		0xb0fb0e4e}
)

const SIZE = 32
const BLOCKSIZE = 64

// -----------------------------------------------------------------------------
// golang/hash/Hash 接口
// -----------------------------------------------------------------------------

// Digest 杂凑过程上下文
type context struct {
	hash   [8]uint32
	buffer [32]byte
}

// New 生成杂凑上下文实例
func New() hash.Hash {
	c := new(context)
	c.Reset()
	return c
}

// Reset 实现 Hash 接口中的 Reset 函数，初始化杂凑上下文
func (c *context) Reset() {
	c.hash = iv
}

// Sum 实现 Hash 接口中的 Sum 函数
func (c *context) Sum(b []byte) []byte {
	mirror := *c
	checkSum := mirror.computeSum()
	return append(b, checkSum[:]...)
}

// Size 实现 Hash 接口中的 Size 函数
func (c *context) Size() int {
	return SIZE
}

// BlockSize 实现 Hash 接口中的 BlockSize 函数
func (c *context) BlockSize() int {
	return BLOCKSIZE
}

// Write 实现 Hash 接口中 io.Writer 接口的 Write 函数
func (c *context) Write(buf []byte) (count int, err error) {
	count = len(buf)
	return
}

// -----------------------------------------------------------------------------
// GB/T 32907 4
// -----------------------------------------------------------------------------

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

// p0
func p0(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 9) ^ bits.RotateLeft32(x, 17)
}

// p1
func p1(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 15) ^ bits.RotateLeft32(x, 23)
}

// -----------------------------------------------------------------------------
// GB/T 32907 5
// -----------------------------------------------------------------------------

// padding GB/T 32907 5.1
func (c *context) padding(message []byte) error {
	return nil
}

func (c *context) computeSum() [SIZE]byte {
	return [SIZE]byte{0}
}
