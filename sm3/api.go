package sm3

import (
	"hash"
)

// -----------------------------------------------------------------------------
// golang/hash/Hash 接口
// -----------------------------------------------------------------------------

// New 生成杂凑上下文实例
func New() hash.Hash {
	c := NewContext()
	c.Reset()
	return c
}

// Reset 实现 Hash 接口中的 Reset 函数，初始化杂凑上下文
func (ctx *Context) Reset() {
	ctx.Init()
}

// Sum 实现 Hash 接口中的 Sum 函数
func (ctx *Context) Sum(inputStream []byte) []byte {
	return ctx.ComputeFromBytes(inputStream).ToBytes()
}

// Size 实现 Hash 接口中的 Size 函数
func (ctx *Context) Size() int {
	return 32
}

// BlockSize 实现 Hash 接口中的 BlockSize 函数
func (ctx *Context) BlockSize() int {
	return 64
}

// Write 实现 Hash 接口中 io.Writer 接口的 Write 函数
func (ctx *Context) Write(newChunk []byte) (int, error) {
	ctx.ComputeFromBytes(newChunk)
	return len(newChunk), nil
}
