package sm2

import (
	"crypto"
	"crypto/elliptic"
	"math/big"

	"github.com/t1anchen/gogmlib/sm3"
	"github.com/t1anchen/gogmlib/utils"
)

// Context 加解密上下文
type Context struct {
	buffer []byte
}

// ToString 字符串表示
func (c *Context) ToString() string {
	return "Hello, World"
}

// -----------------------------------------------------------------------------
// GB/T 32918 3 符号和术语
// -----------------------------------------------------------------------------

// 3.1

type PubKey struct {
	elliptic.Curve
	xA, yA *big.Int
}

type PrivKey struct {
	PubKey
	D *big.Int
}

func (priv *PrivKey) StandardPubKey() crypto.PublicKey {
	return &priv.PubKey
}

var ENTLA, IDA, a, b, xG, yG, xA, yA []byte

// -----------------------------------------------------------------------------
// GB/T 32918 5 数字签名算法
// -----------------------------------------------------------------------------

// KeyPair 5.1.3
type KeyPair struct {
	dA PrivKey
	PA PubKey
}

// H256 5.1.4.2
func H256(msg []byte) []byte {
	return sm3.NewContext().ComputeFromBytes(msg).ToBytes()
}

// ZA 5.1.4.4
var ZA = H256(utils.BytesConcat(ENTLA, IDA, a, b, xG, yG, xA, yA))

// Signature 5.2.1
type Signature struct {
	R, S *big.Int
}

// Sign 5.2.1
func Sign(M []byte) Signature {
	var r, s *big.Int

	// A1
	// Mslide := utils.BytesConcat(ZA, M)

	// A2
	// e := H256(Mslide)

	// A3
	// A4
	// A5
	// A6
	// A7

	return Signature{r, s}
}

// -----------------------------------------------------------------------------
// GB/T 32918 6 密钥交换协议
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// GB/T 32918 7 公开密钥加密
// -----------------------------------------------------------------------------
