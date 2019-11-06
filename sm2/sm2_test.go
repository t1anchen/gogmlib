package sm2

import (
	"crypto/rand"
	"testing"

	"github.com/t1anchen/gogmlib/utils"
)

func TestGenKey(t *testing.T) {
	sk, pk, err := GenKey(rand.Reader)
	t.Logf("pk = %s\nsk = %s\nerr = %s", pk, sk, err)
}

func TestSign(t *testing.T) {
	sk := &PrivKey{
		D:     utils.NewBigIntFromHexString("5DD701828C424B84C5D56770ECF7C4FE882E654CAC53C7CC89A66B1709068B9D"),
		Curve: GetSm2P256()}
	t.Logf("priv.D = %x", sk.D)
	sig, err := sk.Sign(nil, utils.HexStringToBytes("0102030405060708010203040506070801020304050607080102030405060708"))
	if err != nil {
		t.Errorf("TestSign:Sign:err = %s\n", err)
		return
	}
	t.Logf("TestSign:Sign:sig = %x\n", sig)
}

func TestSignToBigInt(t *testing.T) {
	sk := &PrivKey{
		D:     utils.NewBigIntFromHexString("5DD701828C424B84C5D56770ECF7C4FE882E654CAC53C7CC89A66B1709068B9D"),
		Curve: GetSm2P256()}
	msg := utils.HexStringToBytes("0102030405060708010203040506070801020304050607080102030405060708")
	actualR, actualS, err := sk.SignToBigInt(
		sm2SignDefaultUserID,
		msg)
	t.Logf("actualR:%x\nactualS:%x\nerr:%v",
		utils.BigIntToBytes(actualR),
		utils.BigIntToBytes(actualS),
		err)

	pk := &PubKey{
		Curve: GetSm2P256(),
		X:     utils.NewBigIntFromHexString("FF6712D3A7FC0D1B9E01FF471A87EA87525E47C7775039D19304E554DEFE0913"),
		Y:     utils.NewBigIntFromHexString("F632025F692776D4C13470ECA36AC85D560E794E1BCCF53D82C015988E0EB956")}

	actualResult := pk.VerifyFromBigInt(nil, msg, actualR, actualS)
	t.Logf("actualResult:%v", actualResult)

}
