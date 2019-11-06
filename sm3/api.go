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
	return c
}

// Reset 实现 Hash 接口中的 Reset 函数，初始化杂凑上下文
func (ctx *Context) Reset() {
	ctx.Init()
}

// Sum 实现 Hash 接口中的 Sum 函数
func (ctx *Context) Sum(inputStream []byte) []byte {
	digest := ctx
	h := digest.checkSum()
	return append(inputStream, h[:]...)
}

// Size 实现 Hash 接口中的 Size 函数
func (ctx *Context) Size() int {
	return 32
}

// BlockSize 实现 Hash 接口中的 BlockSize 函数
func (ctx *Context) BlockSize() int {
	return BlockSizeInByte
}

// Write 实现 Hash 接口中的 Write 方法（来自 io.Writer）
func (ctx *Context) Write(newChunk []byte) (int, error) {
	return ctx.Update(newChunk)
}
