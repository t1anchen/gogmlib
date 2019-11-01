package sm3

import (
	"hash"
)

// -----------------------------------------------------------------------------
// 应用接口
// -----------------------------------------------------------------------------

// Checksum sm3面向流的 API
func Checksum(message []byte) []byte {
	// TODO
	return nil
}

// -----------------------------------------------------------------------------
// golang/hash/Hash 接口
// -----------------------------------------------------------------------------

// New 生成杂凑上下文实例
func New() hash.Hash {
	c := new(Context)
	c.Reset()
	return c
}

// Reset 实现 Hash 接口中的 Reset 函数，初始化杂凑上下文
func (ctx *Context) Reset() {
	for i, word := range iv {
		ctx.state[i] = word
	}
	ctx.bitCount = 0
	ctx.buffer = [16]uint32{}
	// ctx.m = []byte{}
}

// Sum 实现 Hash 接口中的 Sum 函数
func (ctx *Context) Sum(inputStream []byte) []byte {
	// ctx.Write(inputStream)
	// msg, err := ctx.Padding()
	// if err != nil {
	// 	log.Fatalf("sm3:Sum:Padding 失败, 错误为%v", err)
	// 	return nil
	// }

	// // 最后一个block处理
	// ctx.Compress(msg, len(msg)/ctx.BlockSize())

	// // TODO: 需要使用 bytes 包和 binary 包的相关组件重构
	// requiredLen := ctx.Size()
	// if cap(inputStream)-len(inputStream) < requiredLen {
	// 	newInputStream := make([]byte, len(inputStream), len(inputStream)+requiredLen)
	// 	copy(newInputStream, inputStream)
	// 	inputStream = newInputStream
	// }

	// outputStream := inputStream[len(inputStream) : len(inputStream)+requiredLen]

	// for i := 0; i < 8; i++ {
	// 	binary.BigEndian.PutUint32(outputStream[i*4:], ctx.state[i])
	// }
	// return outputStream
	return nil
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
	// newChunkLen := len(newChunk)
	// ctx.bitCount += uint64(len(newChunk) * 8)

	// msg := append(ctx.m, newChunk...)
	// n := len(msg) / ctx.BlockSize() // blocks
	// ctx.Compress(msg, n)

	// ctx.m = msg[n*ctx.BlockSize():]

	return 0, nil
}
