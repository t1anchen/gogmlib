package sm3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"

	"github.com/t1anchen/gogmlib/utils"
)

const (
	BlockSizeInByte  = 16
	DigestSizeInByte = 32
)

// -----------------------------------------------------------------------------
// 数据结构
// -----------------------------------------------------------------------------

// Context 杂凑过程上下文
type Context struct {
	state     [8]uint32
	buffer    [16]uint32
	w         [68]uint32
	wp        [64]uint32
	xBuf      [4]byte
	xBufOff   int32
	xOff      int32
	byteCount int64
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
	}
	return 0x79cc4519
}

// ff 4.3
func ff(j int, x, y, z uint32) uint32 {
	if j >= 16 {
		return ((x | y) & (x | z) & (y | z))
	}
	return x ^ y ^ z
}

// gg 4.3
func gg(j int, x, y, z uint32) uint32 {
	if j >= 16 {
		return ((x & y) | ((^x) & z))
	}
	return x ^ y ^ z
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

// NewContext 5.1
func NewContext() *Context {
	var ctx Context
	ctx.Init()
	return &ctx
}

// Init 5.1
func (ctx *Context) Init() *Context {
	ctx.byteCount = 0
	ctx.xBufOff = 0

	for i := 0; i < len(ctx.xBuf); i++ {
		ctx.xBuf[i] = 0
	}
	for i := 0; i < 8; i++ {
		ctx.state[i] = iv[i]
	}
	for i := 0; i < len(ctx.w); i++ {
		ctx.w[i] = 0
	}
	ctx.xOff = 0
	return ctx
}

// Padding 5.2 填充
func Padding(message []byte) []byte {
	msgLen := len(message)
	msgBuf := bytes.NewBuffer(message)
	msgBuf.WriteByte(0x80)
	for msgBuf.Len()%64 != 56 {
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

// CompressFunction 5.3.3 压缩函数
func (ctx *Context) CompressFunction(w *[68]uint32, wp *[64]uint32) {
	var ss1, ss2, tt1, tt2 uint32

	a, b, c, d, e, f, g, h :=
		ctx.state[0],
		ctx.state[1],
		ctx.state[2],
		ctx.state[3],
		ctx.state[4],
		ctx.state[5],
		ctx.state[6],
		ctx.state[7]

	for j := 0; j < 64; j++ {
		ss1 = rotl32(rotl32(a, 12)+e+rotl32(t(j), j), 7)
		ss2 = ss1 ^ rotl32(a, 12)
		tt1 = ff(j, a, b, c) + d + ss2 + wp[j]
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
	ctx.state[0] ^= a
	ctx.state[1] ^= b
	ctx.state[2] ^= c
	ctx.state[3] ^= d
	ctx.state[4] ^= e
	ctx.state[5] ^= f
	ctx.state[6] ^= g
	ctx.state[7] ^= h
}

// ComputeFromBytes 总过程
func (ctx *Context) ComputeFromBytes(message []byte) *Context {
	padded := bytes.NewBuffer(Padding(message))
	blockBuf := bytes.NewBuffer(padded.Next(64))
	for blockBuf.Len() > 0 {
		binary.Read(blockBuf, binary.BigEndian, &ctx.buffer)
		ctx.MessageExpansion(&ctx.w, &ctx.wp)
		ctx.CompressFunction(&ctx.w, &ctx.wp)
		blockBuf = bytes.NewBuffer(padded.Next(64))
	}
	return ctx
}

func (ctx *Context) processBlock() *Context {
	ctx.MessageExpansion(&ctx.w, &ctx.wp)
	ctx.CompressFunction(&ctx.w, &ctx.wp)
	ctx.xOff = 0
	return ctx
}

func (ctx *Context) Update(payload []byte) (n int, err error) {
	payloadLen := len(payload)
	i := 0
	if ctx.xBufOff != 0 {
		for i < payloadLen {
			ctx.xBuf[ctx.xBufOff] = payload[i]
			ctx.xBufOff++
			i++
			if ctx.xBufOff == 4 {
				ctx.processWord(ctx.xBuf[:], 0)
				ctx.xBufOff = 0
				break
			}
		}
	}

	limit := ((payloadLen - i) & ^3) + i
	for ; i < limit; i += 4 {
		ctx.processWord(payload, int32(i))
	}

	for i < payloadLen {
		ctx.xBuf[ctx.xBufOff] = payload[i]
		ctx.xBufOff++
		i++
	}

	ctx.byteCount += int64(payloadLen)

	n = payloadLen
	return
}

func (ctx *Context) processWord(wordBuf []byte, offset int32) {
	n := binary.BigEndian.Uint32(wordBuf[offset : offset+4])

	ctx.buffer[ctx.xOff] = n
	ctx.xOff++

	if ctx.xOff >= 16 {
		ctx.processBlock()
	}
}

func (ctx *Context) processLength(bitLength int64) {
	if ctx.xOff > (BlockSizeInByte - 2) {
		ctx.buffer[ctx.xOff] = 0
		ctx.xOff++

		ctx.processBlock()
	}

	for ; ctx.xOff < (BlockSizeInByte - 2); ctx.xOff++ {
		ctx.buffer[ctx.xOff] = 0
	}

	ctx.buffer[ctx.xOff] = uint32(bitLength >> 32)
	ctx.xOff++
	ctx.buffer[ctx.xOff] = uint32(bitLength)
	ctx.xOff++
}

func (ctx *Context) finish() {
	bitLength := ctx.byteCount << 3

	ctx.Write([]byte{128})

	for ctx.xBufOff != 0 {
		ctx.Write([]byte{0})
	}

	ctx.processLength(bitLength)

	ctx.processBlock()
}

func (ctx *Context) checkSum() [32]byte {
	ctx.finish()
	vlen := len(ctx.state)
	var out [32]byte
	for i := 0; i < vlen; i++ {
		binary.BigEndian.PutUint32(out[i*4:(i+1)*4], ctx.state[i])
	}
	return out
}

// ComputeFromString 以 string 输入
func (ctx *Context) ComputeFromString(s string) *Context {
	return ctx.ComputeFromBytes([]byte(s))
}

// ComputeFromWords 以 Word 输入
func (ctx *Context) ComputeFromWords(words []uint32) *Context {
	return ctx.ComputeFromBytes(utils.WordsToBytes(words))
}

// ToWords 以 [8]uint32 输出
func (ctx *Context) ToWords() []uint32 {
	return ctx.state[:]
}

// ToBytes 以 Bytes 输出
func (ctx *Context) ToBytes() []byte {
	return utils.WordsToBytes(ctx.state[:])
}

// ToHexString 以 string 输出
func (ctx *Context) ToHexString() string {
	return fmt.Sprintf("%x", ctx.ToBytes())
}
