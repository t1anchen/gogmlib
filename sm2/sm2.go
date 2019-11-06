package sm2

import (
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"

	"github.com/t1anchen/gogmlib/utils"
)

const (
	BitSize      = 256
	KeyBytes     = (BitSize + 7) / 8
	UnCompressed = 0x04
)

// PubKey 公钥
type PubKey struct {
	X, Y  *big.Int
	Curve Curve
}

// PrivKey 私钥
type PrivKey struct {
	D     *big.Int
	Curve Curve
}

// Signature 签名
type Signature struct {
	R, S *big.Int
}

// Cipher 密文
type Cipher struct {
	X, Y *big.Int
	C3   []byte
	C2   []byte
}

// ToUncompressedBytes 公钥未压缩字节流
func (pk *PubKey) ToUncompressedBytes() []byte {
	xBytes := utils.BigIntToBytes(pk.X)
	yBytes := utils.BigIntToBytes(pk.Y)
	xLen := len(xBytes)
	yLen := len(yBytes)

	pkBuf := make([]byte, 1+KeyBytes*2)
	pkBuf[0] = UnCompressed
	if xLen > KeyBytes {
		copy(pkBuf[1:1+KeyBytes], xBytes[xLen-KeyBytes:])
	} else if xLen < KeyBytes {
		copy(pkBuf[1+(KeyBytes-xLen):1+KeyBytes], xBytes)
	} else {
		copy(pkBuf[1:1+KeyBytes], xBytes)
	}

	if yLen > KeyBytes {
		copy(pkBuf[1+KeyBytes:], yBytes[yLen-KeyBytes:])
	} else if yLen < KeyBytes {
		copy(pkBuf[1+KeyBytes+(KeyBytes-yLen):], yBytes)
	} else {
		copy(pkBuf[1+KeyBytes:], yBytes)
	}
	return pkBuf
}

// ToBytes 公钥的字节流
func (pk *PubKey) ToBytes() []byte {
	return pk.ToUncompressedBytes()[1:]
}

// ToBytes 私钥的字节流
func (sk *PrivKey) ToBytes() []byte {
	dBytes := utils.BigIntToBytes(sk.D)
	dLen := len(dBytes)

	if dLen > KeyBytes {
		skBuf := make([]byte, KeyBytes)
		copy(skBuf, dBytes[dLen-KeyBytes:])
		return skBuf
	} else if dLen < KeyBytes {
		skBuf := make([]byte, KeyBytes)
		copy(skBuf[KeyBytes-dLen:], dBytes)
		return skBuf
	} else {
		return dBytes
	}
}

// -----------------------------------------------------------------------------
// ietf/draft-shen-sm2-ecdsa-02 5 数字签名算法
// -----------------------------------------------------------------------------
var sm2P256 Curve

// GenKey 生成密钥对
func GenKey(rand io.Reader) (*PrivKey, *PubKey, error) {
	privFromGoCrypto, xFromGoCrypto, yFromGoCrypto, err := elliptic.GenerateKey(sm2P256, rand)
	if err != nil {
		return nil, nil, err
	}
	privKey := &PrivKey{
		D:     utils.NewBigIntFromBytes(privFromGoCrypto),
		Curve: sm2P256}
	pubKey := &PubKey{
		Curve: sm2P256,
		X:     xFromGoCrypto,
		Y:     yFromGoCrypto}
	return privKey, pubKey, nil
}

// GenPubKey 生成公钥
func (sk *PrivKey) GenPubKey() *PubKey {
	pk := new(PubKey)
	pk.Curve = sk.Curve
	pk.X, pk.Y = sk.Curve.ScalarBaseMult(utils.BigIntToBytes(sk.D))
	return pk
}

func pickK(randomStream io.Reader, upper *big.Int) (*big.Int, error) {
	one := utils.NewBigIntFromOne()
	var k *big.Int
	var err error
	for {
		k, err = rand.Int(randomStream, upper)
		if err != nil {
			return nil, err
		}
		if k.Cmp(one) >= 0 {
			return k, err
		}
	}
}

// -----------------------------------------------------------------------------
// ietf/draft-shen-sm2-ecdsa-02 6 密钥交换协议
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// ietf/draft-shen-sm2-ecdsa-02 7 公开密钥加密
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// golang package 本身
// -----------------------------------------------------------------------------
func init() {
	sm2P256 = *InitWithRecommendedParams()
}
