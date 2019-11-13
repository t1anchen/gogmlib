package sm2

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"

	"github.com/t1anchen/gogmlib/sm3"
	"github.com/t1anchen/gogmlib/utils"
)

const (
	BitSize      = 256
	KeyBytes     = (BitSize + 7) / 8
	UnCompressed = 0x04
)

// PubKey 公钥
type PubKey struct {
	X     *big.Int `json:"x"`
	Y     *big.Int `json:"y"`
	Curve Curve    `json:"curve"`
}

func (pk *PubKey) String() string {
	b, _ := json.Marshal(pk)
	return fmt.Sprintf("%s", b)
}

// PrivKey 私钥
type PrivKey struct {
	D     *big.Int `json:"d"`
	Curve Curve    `json:"curve"`
}

func (sk *PrivKey) String() string {
	b, _ := json.Marshal(sk)
	return fmt.Sprintf("%s", b)
}

// Signature 签名
type Signature struct {
	R, S *big.Int
}

// ASN1DERCipher 密文
type ASN1DERCipher struct {
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

var (
	sm2H                 = utils.NewBigIntFromOne()
	sm2SignDefaultUserID = utils.HexStringToBytes("31323334353637383132333435363738")
)

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

func getZ(hashProvider hash.Hash, curve *Curve, pubX *big.Int, pubY *big.Int, userID []byte) []byte {
	hashProvider.Reset()

	userIDLen := uint16(len(userID) * 8)
	var userIDLenBytes [2]byte
	binary.BigEndian.PutUint16(userIDLenBytes[:], userIDLen)
	hashProvider.Write(userIDLenBytes[:])
	if userID != nil && len(userID) > 0 {
		hashProvider.Write(userID)
	}

	hashProvider.Write(utils.BigIntToBytes(curve.A))
	hashProvider.Write(utils.BigIntToBytes(curve.B))
	hashProvider.Write(utils.BigIntToBytes(curve.Gx))
	hashProvider.Write(utils.BigIntToBytes(curve.Gy))
	hashProvider.Write(utils.BigIntToBytes(pubX))
	hashProvider.Write(utils.BigIntToBytes(pubY))
	return hashProvider.Sum(nil)
}

func genE(hashProvider hash.Hash, curve *Curve, pubX *big.Int, pubY *big.Int, userId []byte, src []byte) *big.Int {
	z := getZ(hashProvider, curve, pubX, pubY, userId)

	hashProvider.Reset()
	hashProvider.Write(z)
	hashProvider.Write(src)
	eHash := hashProvider.Sum(nil)
	return utils.NewBigIntFromBytes(eHash)
}

// SignToASN1DER 生成签名并输出成 ASN.1 DER 格式
func (sk *PrivKey) SignToASN1DER(userID []byte, msg []byte) ([]byte, error) {
	r, s, signErr := sk.SignToBigInt(userID, msg)
	if signErr != nil {
		return nil, signErr
	}

	encoded, encErr := asn1.Marshal(Signature{r, s})
	if encErr != nil {
		return nil, encErr
	}

	return encoded, nil

}

// SignToBigInt 生成签名并输出成大数
func (sk *PrivKey) SignToBigInt(userID []byte, msg []byte) (r, s *big.Int, err error) {
	digest := sm3.New()
	pubX, pubY := sk.Curve.ScalarBaseMult(utils.BigIntToBytes(sk.D))
	if userID == nil {
		userID = sm2SignDefaultUserID
	}
	fmt.Printf("userID:%x\nmsg:%x\npubX:%x\npubY:%x\n",
		userID, msg, pubX.Bytes(), pubY.Bytes())
	e := genE(digest, &sk.Curve, pubX, pubY, userID, msg)
	fmt.Printf("e:%x\n", e.Bytes())

	bigIntZero := utils.NewBigIntFromZero()
	bigIntOne := utils.NewBigIntFromOne()
	for {
		var k *big.Int
		var err error
		for {
			k, err = pickK(rand.Reader, sk.Curve.N)
			if err != nil {
				return nil, nil, err
			}
			px, _ := sk.Curve.ScalarBaseMult(utils.BigIntToBytes(k))
			r = utils.BigIntAdd(e, px)
			r = utils.BigIntMod(r, sk.Curve.N)

			rk := new(big.Int).Set(r)
			rk = rk.Add(rk, k)
			if r.Cmp(bigIntZero) != 0 && rk.Cmp(sk.Curve.N) != 0 {
				break
			}
		}

		dPlus1ModN := utils.BigIntAdd(sk.D, bigIntOne)
		dPlus1ModN = utils.BigIntModInverse(dPlus1ModN, sk.Curve.N)
		s = utils.BigIntMul(r, sk.D)
		s = utils.BigIntSub(k, s)
		s = utils.BigIntMod(s, sk.Curve.N)
		s = utils.BigIntMul(dPlus1ModN, s)
		s = utils.BigIntMod(s, sk.Curve.N)

		if s.Cmp(bigIntZero) != 0 {
			break
		}
	}

	return r, s, nil
}

func (pk *PubKey) VerifyFromBigInt(userID []byte, msg []byte, r, s *big.Int) bool {
	bigIntOne := utils.NewBigIntFromOne()
	if r.Cmp(bigIntOne) == -1 || r.Cmp(pk.Curve.N) >= 0 {
		return false
	}
	if s.Cmp(bigIntOne) == -1 || s.Cmp(pk.Curve.N) >= 0 {
		return false
	}

	digest := sm3.New()
	if userID == nil {
		userID = sm2SignDefaultUserID
	}
	e := genE(digest, &pk.Curve, pk.X, pk.Y, userID, msg)

	intZero := utils.NewBigIntFromZero()
	t := utils.BigIntAdd(r, s)
	t = utils.BigIntMod(t, pk.Curve.N)
	if t.Cmp(intZero) == 0 {
		return false
	}

	sgx, sgy := pk.Curve.ScalarBaseMult(s.Bytes())
	tpx, tpy := pk.Curve.ScalarMult(pk.X, pk.Y, t.Bytes())
	x, y := pk.Curve.Add(sgx, sgy, tpx, tpy)
	if IsPointInfinity(x, y) {
		return false
	}

	expectedR := utils.BigIntAdd(e, x)
	expectedR = utils.BigIntMod(expectedR, pk.Curve.N)
	return expectedR.Cmp(r) == 0
}

// VerifyFromASN1DER 验证签名
func (pk *PubKey) VerifyFromASN1DER(userID []byte, msg []byte, sigBytes []byte) bool {
	sig := new(Signature)
	_, err := asn1.Unmarshal(sigBytes, sig)
	if err != nil {
		fmt.Errorf("%s", err)
		return false
	}
	return pk.VerifyFromBigInt(userID, msg, sig.R, sig.S)
}

// -----------------------------------------------------------------------------
// ietf/draft-shen-sm2-ecdsa-02 6 密钥交换协议
// -----------------------------------------------------------------------------

func kdf(hashProvider hash.Hash, zX, zY *big.Int, msg []byte) {
	bufSize := 4
	if bufSize < hashProvider.BlockSize() {
		bufSize = hashProvider.BlockSize()
	}
	buf := make([]byte, bufSize)

	msgLen := len(msg)
	zXBytes := zX.Bytes()
	zYBytes := zY.Bytes()
	offset := 0
	roundCount := uint32(0)
	for offset < msgLen {
		hashProvider.Reset()
		hashProvider.Write(zXBytes)
		hashProvider.Write(zYBytes)
		roundCount++
		binary.BigEndian.PutUint32(buf, roundCount)
		hashProvider.Write(buf[:4])
		tmp := hashProvider.Sum(nil)
		copy(buf[:bufSize], tmp[:bufSize])

		xorLen := msgLen - offset
		if xorLen > hashProvider.BlockSize() {
			xorLen = hashProvider.BlockSize()
		}
		utils.BytesXor(msg[offset:], buf, xorLen)
		offset += xorLen
	}
}

// -----------------------------------------------------------------------------
// ietf/draft-shen-sm2-ecdsa-02 7 公开密钥加密
// -----------------------------------------------------------------------------

func isNotEncrypted(encData []byte, in []byte) bool {
	encDataLen := len(encData)
	for i := 0; i != encDataLen; i++ {
		if encData[i] != in[i] {
			return false
		}
	}
	return true
}

func (pk *PubKey) Encrypt(msg []byte) ([]byte, error) {
	c2 := make([]byte, len(msg))
	copy(c2, msg)
	var c1 []byte
	digest := sm3.New()
	var kPBx, kPBy *big.Int
	for {
		k, err := pickK(rand.Reader, pk.Curve.N)
		if err != nil {
			return nil, err
		}
		kBytes := k.Bytes()
		c1x, c1y := pk.Curve.ScalarBaseMult(kBytes)
		c1 = elliptic.Marshal(pk.Curve, c1x, c1y)
		kPBx, kPBy = pk.Curve.ScalarMult(pk.X, pk.Y, kBytes)
		kdf(digest, kPBx, kPBy, c2)

		if !isNotEncrypted(c2, msg) {
			break
		}
	}

	digest.Reset()
	digest.Write(kPBx.Bytes())
	digest.Write(msg)
	digest.Write(kPBy.Bytes())
	c3 := digest.Sum(nil)

	c1Len := len(c1)
	c2Len := len(c2)
	c3Len := len(c3)
	result := make([]byte, c1Len+c2Len+c3Len)
	copy(result[:c1Len], c1)
	copy(result[c1Len:c1Len+c2Len], c2)
	copy(result[c1Len+c2Len:], c3)
	return result, nil

}

func (sk *PrivKey) Decrypt(in []byte) ([]byte, error) {
	c1Len := ((sk.Curve.BitSize+7)/8)*2 + 1
	c1 := make([]byte, c1Len)
	copy(c1, in[:c1Len])
	c1x, c1y := elliptic.Unmarshal(sk.Curve, c1)
	sx, sy := sk.Curve.ScalarMult(c1x, c1y, sm2H.Bytes())
	if IsPointInfinity(sx, sy) {
		return nil, errors.New("[h]C1 at infinity")
	}
	c1x, c1y = sk.Curve.ScalarMult(c1x, c1y, sk.D.Bytes())

	digest := sm3.New()
	c3Len := digest.Size()
	c2Len := len(in) - c1Len - c3Len
	c2 := make([]byte, c2Len)
	copy(c2, in[c1Len:c1Len+c2Len])
	kdf(digest, c1x, c1y, c2)

	digest.Reset()
	digest.Write(c1x.Bytes())
	digest.Write(c2)
	digest.Write(c1y.Bytes())
	c3 := digest.Sum(nil)

	if !bytes.Equal(c3, in[c1Len+c2Len:]) {
		return nil, errors.New("invalid cipher text")
	}
	return c2, nil
}

func MarshalEncryptedToASN1DER(in []byte) ([]byte, error) {
	paramLen := utils.BitsLenToBytesLen(sm2P256.Params().BitSize)

	c1x := make([]byte, paramLen)
	c1y := make([]byte, paramLen)
	c2Len := len(in) - (1 + paramLen*2) - sm3.DigestSizeInByte
	c2 := make([]byte, c2Len)
	c3 := make([]byte, sm3.DigestSizeInByte)
	pos := 1

	copy(c1x, in[pos:pos+paramLen])
	pos += paramLen

	copy(c1y, in[pos:pos+paramLen])
	pos += paramLen

	copy(c2, in[pos:pos+c2Len])
	pos += c2Len

	copy(c3, in[pos:pos+sm3.DigestSizeInByte])

	nc1x := utils.NewBigIntFromBytes(c1x)
	nc1y := utils.NewBigIntFromBytes(c1y)
	result, err := asn1.Marshal(ASN1DERCipher{nc1x, nc1y, c3, c2})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UnmarshalEncryptedFromASN1DER(in []byte) ([]byte, error) {
	encrypted := new(ASN1DERCipher)
	_, err := asn1.Unmarshal(in, encrypted)
	if err != nil {
		return nil, err
	}

	c1x := utils.BigIntToBytes(encrypted.X)
	c1y := utils.BigIntToBytes(encrypted.Y)
	c1xLen := len(c1x)
	c1yLen := len(c1y)
	c2Len := len(encrypted.C2)
	c3Len := len(encrypted.C3)
	result := make([]byte, 1+c1xLen+c1yLen+c2Len+c3Len)
	pos := 0

	result[pos] = UnCompressed
	pos++

	copy(result[pos:pos+c1xLen], c1x)
	pos += c1xLen

	copy(result[pos:pos+c1yLen], c1y)
	pos += c1yLen

	copy(result[pos:pos+c2Len], encrypted.C2)
	pos += c2Len

	copy(result[pos:pos+c3Len], encrypted.C3)

	return result, nil
}

// -----------------------------------------------------------------------------
// golang package 本身
// -----------------------------------------------------------------------------
func init() {
	sm2P256 = *InitWithRecommendedParams()
}
