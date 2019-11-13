package sm2

import (
	"crypto/elliptic"
	"math/big"

	"github.com/t1anchen/gogmlib/utils"
)

// Curve 曲线和曲线参数 3.1
type Curve struct {
	*elliptic.CurveParams
	RInverse *big.Int
	A        *big.Int
}

var sm2P256 Curve

// Init 初始化曲线
func Init(A *big.Int, RInverse *big.Int, givenCurveParams *elliptic.CurveParams) *Curve {
	var sm2P256 Curve
	sm2P256.A = A
	sm2P256.CurveParams = givenCurveParams
	sm2P256.RInverse = RInverse
	return &sm2P256
}

// InitWithRecommendedParams 使用推荐参数初始化曲线
func InitWithRecommendedParams() *Curve {
	return Init(
		utils.NewBigIntFromHexString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFC"),
		utils.NewBigIntFromHexString("7ffffffd80000002fffffffe000000017ffffffe800000037ffffffc80000002"),
		&elliptic.CurveParams{
			Name:    "SM2-P-256",
			P:       utils.NewBigIntFromHexString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF"),
			N:       utils.NewBigIntFromHexString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123"),
			B:       utils.NewBigIntFromHexString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93"),
			Gx:      utils.NewBigIntFromHexString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7"),
			Gy:      utils.NewBigIntFromHexString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0"),
			BitSize: 256})
}

func init() {
	sm2P256 = *InitWithRecommendedParams()
}

func GetSm2P256() Curve {
	return sm2P256
}

func IsPointInfinity(x, y *big.Int) bool {
	if x.Sign() == 0 && y.Sign() == 0 {
		return true
	}
	return false
}
